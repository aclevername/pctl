package profile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/otiai10/copy"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/weaveworks/pctl/pkg/git"
	"github.com/weaveworks/pctl/pkg/profile/artifact"
	"github.com/weaveworks/pctl/pkg/profile/builder"
)

// ArtifactsMaker can create a list of artifacts.
//go:generate counterfeiter -o fakes/artifacts_maker.go . ArtifactsMaker
type ArtifactsMaker interface {
	Make(installation profilesv1.ProfileInstallation) error
}

// MakerConfig contains all configuration properties for the Artifacts Maker.
type MakerConfig struct {
	GitClient        git.Git
	RootDir          string
	GitRepoNamespace string
	GitRepoName      string
	ProfileName      string
}

// ProfilesArtifactsMaker creates a list of artifacts from profiles data.
type ProfilesArtifactsMaker struct {
	MakerConfig

	Builder      builder.Builder
	nestedName   string
	profileRepos []string
	cloneCache   map[string]string
}

// NewProfilesArtifactsMaker creates a new profiles artifacts maker.
func NewProfilesArtifactsMaker(cfg MakerConfig) *ProfilesArtifactsMaker {
	return &ProfilesArtifactsMaker{
		MakerConfig: cfg,
		Builder: &builder.ArtifactBuilder{
			Config: builder.Config{
				GitRepositoryName:      cfg.GitRepoName,
				GitRepositoryNamespace: cfg.GitRepoNamespace,
				RootDir:                cfg.RootDir,
			},
		},
		cloneCache: make(map[string]string),
	}
}

// Make generates artifacts without owners for manual applying to a personal cluster.
func (pa *ProfilesArtifactsMaker) Make(installation profilesv1.ProfileInstallation) error {
	artifacts, err := profilesArtifactsMaker(pa, installation)
	if err != nil {
		return fmt.Errorf("failed to build artifact: %w", err)
	}
	profileRootdir := filepath.Join(pa.RootDir, pa.ProfileName)
	artifactsRootDir := filepath.Join(profileRootdir, "artifacts")
	defer pa.cleanCloneCache()
	for _, artifact := range artifacts {
		artifactDir := filepath.Join(artifactsRootDir, artifact.Name)
		if err := os.MkdirAll(artifactDir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create directory")
		}
		if artifact.RepoURL != "" {
			if err := pa.getRepositoryLocalArtifacts(artifact, artifactDir); err != nil {
				return fmt.Errorf("failed to get package local artifacts: %w", err)
			}
		}
		for _, obj := range artifact.Objects {
			name := obj.GetObjectKind().GroupVersionKind().Kind
			if obj.Name != "" {
				name = obj.Name
			}
			filename := filepath.Join(artifactDir, fmt.Sprintf("%s.%s", name, "yaml"))
			if obj.Path != "" {
				subFolder := filepath.Join(artifactDir, obj.Path)
				if err := os.MkdirAll(subFolder, 0755); err != nil && !os.IsExist(err) {
					return fmt.Errorf("failed to create directory")
				}
				filename = filepath.Join(subFolder, fmt.Sprintf("%s.%s", name, "yaml"))
			}
			if err := pa.generateOutput(filename, obj.Object); err != nil {
				return err
			}
		}
		// if we have a local resource, write out the kustomization yaml limiting its visibility.
		if artifact.Kustomize.LocalResourceLimiter != nil {
			// This is helmRelease related so it must be inside the sub-folder for the helm release.
			filename := filepath.Join(artifactDir, artifact.SubFolder, "kustomization.yaml")
			if err := writeOutKustomizeResource(artifact.Kustomize.LocalResourceLimiter, filename); err != nil {
				return err
			}
		}
		filename := filepath.Join(artifactDir, "kustomization.yaml")
		if err := writeOutKustomizeResource(artifact.Kustomize.ObjectWrapper, filename); err != nil {
			return err
		}
	}
	return pa.generateOutput(filepath.Join(profileRootdir, "profile-installation.yaml"), &installation)
}

// writeOutKustomizeResource writes out kustomization resource data if set to a specific file.
func writeOutKustomizeResource(kustomize *types.Kustomization, filename string) error {
	data, err := yaml.Marshal(kustomize)
	if err != nil {
		return fmt.Errorf("failed to marshal kustomize resource: %w", err)
	}
	if err = os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

// getRepositoryLocalArtifacts clones all repository local artifacts so they can be copied over to the flux repository.
func (pa *ProfilesArtifactsMaker) getRepositoryLocalArtifacts(a artifact.Artifact, artifactDir string) error {
	var (
		tmp string
		err error
	)
	if v, ok := pa.cloneCache[cloneCacheKey(a.RepoURL, a.Branch)]; ok {
		tmp = v
	} else {
		u := uuid.NewString()[:6]
		tmp, err = ioutil.TempDir("", "clone_git_repo_"+u)
		if err != nil {
			return fmt.Errorf("failed to create temp folder: %w", err)
		}
		if err := pa.GitClient.Clone(a.RepoURL, a.Branch, tmp); err != nil {
			return fmt.Errorf("failed to sparse clone folder with url: %s; branch: %s; with error: %w",
				a.RepoURL,
				a.Branch,
				err)
		}
		pa.cloneCache[cloneCacheKey(a.RepoURL, a.Branch)] = tmp
	}

	for _, path := range a.PathsToCopy {
		// nginx/chart/...
		if strings.Contains(path, string(os.PathSeparator)) {
			path = filepath.Dir(path)
		}
		fullPath := filepath.Join(tmp, a.SparseFolder, path)
		dest := filepath.Join(artifactDir, path)
		if a.SubFolder != "" {
			dest = filepath.Join(artifactDir, a.SubFolder, path)
		}
		if err := copy.Copy(fullPath, dest); err != nil {
			return fmt.Errorf("failed to move folder: %w", err)
		}
	}
	return nil
}

// cleanCloneCache clears all cached cloned folders if there are any.
func (pa *ProfilesArtifactsMaker) cleanCloneCache() {
	for _, c := range pa.cloneCache {
		if err := os.RemoveAll(c); err != nil {
			fmt.Printf("failed to remove %s cache, please clean by hand", c)
		}
	}
}

func (pa *ProfilesArtifactsMaker) generateOutput(filename string, o runtime.Object) error {
	e := kjson.NewSerializerWithOptions(kjson.DefaultMetaFactory, nil, nil, kjson.SerializerOptions{Yaml: true, Strict: true})
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Printf("Failed to properly close file %s\n", f.Name())
		}
	}(f)
	if err := e.Encode(o, f); err != nil {
		return err
	}
	return nil

}

func profileRepo(installation profilesv1.ProfileInstallation) string {
	if installation.Spec.Source.Tag != "" {
		return installation.Spec.Source.URL + ":" + installation.Spec.Source.Tag
	}
	return installation.Spec.Source.URL + ":" + installation.Spec.Source.Branch + ":" + installation.Spec.Source.Path
}

func containsKey(list []string, key string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

func cloneCacheKey(url, branch string) string {
	return fmt.Sprintf("%s:%s", url, branch)
}
