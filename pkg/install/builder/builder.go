package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/dependency"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/otiai10/copy"
	"github.com/weaveworks/pctl/pkg/install/artifact"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/kustomize/api/types"
)

const (
	defaultValuesKey           = "default-values.yaml"
	helmChartLocation          = "helm-chart"
	kustomizeWrapperObjectName = "kustomize-flux.yaml"
)

// Builder can build an artifacts from an installation and a profile artifact.
//go:generate counterfeiter -o fakes/fake_builder.go . Builder
type Builder interface {
	Write(installation profilesv1.ProfileInstallation, artifacts []artifact.Artifact, repoLocation map[string]string) error
}

var _ Builder = &ArtifactBuilder2{}

// ArtifactBuilder will build helm chart resources.
type ArtifactBuilder2 struct {
	GitRepositoryName      string
	GitRepositoryNamespace string
	RootDir                string
}

// Build a single artifact from a profile artifact and installation.
func (c *ArtifactBuilder2) Write(installation profilesv1.ProfileInstallation, artifacts []artifact.Artifact, repoLocation map[string]string) error {
	for _, a := range artifacts {
		var deps []artifact.Artifact
		for _, dep := range a.DependsOn {
			d, ok := containsArtifact(artifacts, dep.Name)
			if !ok {
				return fmt.Errorf("%s's depending artifact %s not found in the list of artifacts", a.Name, dep.Name)
			}
			deps = append(deps, d)
		}
		if err := validateArtifact(a.Artifact); err != nil {
			return fmt.Errorf("invalid artifact: %w", err)
		}

		if a.Chart != nil {
			if err := c.writeChartArtifact(installation, a, deps, repoLocation); err != nil {
				return err
			}
		} else if a.Kustomize != nil {
			if err := c.writeKustomizeArtifact(installation, a, deps, repoLocation); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no artifact type set")
		}
	}
	return c.writeResourceWithName(&installation, filepath.Join(c.RootDir, "profile-installation.yaml"))
}

func validateArtifact(a profilesv1.Artifact) error {
	if a.Chart != nil {
		if a.Profile != nil {
			return apis.ErrMultipleOneOf("chart", "profile")
		}
		if a.Kustomize != nil {
			return apis.ErrMultipleOneOf("chart", "kustomize")
		}
		if a.Chart.Path != "" && a.Chart.URL != "" {
			return apis.ErrMultipleOneOf("chart.path", "chart.url")
		}
	}
	if a.Kustomize != nil && a.Profile != nil {
		return apis.ErrMultipleOneOf("kustomize", "profile")
	}
	return nil
}

func (c *ArtifactBuilder2) writeKustomizeArtifact(installation profilesv1.ProfileInstallation, a artifact.Artifact, deps []artifact.Artifact, repoLocations map[string]string) error {
	//TODO validate kustomize
	if c.GitRepositoryNamespace == "" || c.GitRepositoryName == "" {
		return fmt.Errorf("in case of local resources, the flux gitrepository object's details must be provided")
	}

	artifactDir := filepath.Join(c.RootDir, "artifacts", a.ParentProfileArtifactName, a.Name)
	if err := c.copyArtifacts(a, a.Kustomize.Path, filepath.Join(artifactDir, a.Kustomize.Path), repoLocations); err != nil {
		return err
	}
	if err := c.writeOutKustomizeResource([]string{kustomizeWrapperObjectName}, artifactDir); err != nil {
		return err
	}

	// wrapper := c.makeKustomization(artifact, path, installation, definitionName, dependencies)
	return c.writeResourceWithName(c.makeKustomization(a, filepath.Join(artifactDir, a.Kustomize.Path), installation, a.ProfilePath, deps), filepath.Join(artifactDir, kustomizeWrapperObjectName))
}

func (c *ArtifactBuilder2) writeChartArtifact(installation profilesv1.ProfileInstallation, a artifact.Artifact, deps []artifact.Artifact, repoLocations map[string]string) error {
	//TODO validate chart
	artifactDir := filepath.Join(c.RootDir, "artifacts", a.ParentProfileArtifactName, a.Name)
	helmChartDir := filepath.Join(artifactDir, helmChartLocation)
	if err := os.MkdirAll(helmChartDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory %w", err)
	}
	var objs []runtime.Object
	helmRelease, cfgMap := c.makeHelmReleaseObjects(a.Artifact, installation, a.ProfilePath)
	if cfgMap != nil {
		objs = append(objs, cfgMap)
	}
	objs = append(objs, helmRelease)
	if a.Chart.Path != "" {
		if c.GitRepositoryNamespace == "" || c.GitRepositoryName == "" {
			return fmt.Errorf("in case of local resources, the flux gitrepository object's details must be provided")
		}
		helmRelease.Spec.Chart.Spec.Chart = filepath.Join(helmChartDir, a.Chart.Path)
		if err := c.copyArtifacts(a, a.Chart.Path, filepath.Join(helmChartDir, a.Chart.Path), repoLocations); err != nil {
			return err
		}

		if err := c.writeOutKustomizeResource([]string{"HelmRelease.yaml"}, helmChartDir); err != nil {
			return err
		}

	} else {
		objs = append(objs, c.makeHelmRepository(a.Chart.URL, a.Chart.Name, installation))
	}

	for _, obj := range objs {
		if err := c.writeResource(obj, helmChartDir); err != nil {
			return err
		}
	}

	if err := c.writeOutKustomizeResource([]string{kustomizeWrapperObjectName}, artifactDir); err != nil {
		return err
	}

	return c.writeResourceWithName(c.makeKustomizeHelmReleaseWrapper(a, installation, a.ProfilePath, helmChartDir, deps), filepath.Join(artifactDir, kustomizeWrapperObjectName))
}

func (c *ArtifactBuilder2) writeOutKustomizeResource(resources []string, dir string) error {
	kustomize := &types.Kustomization{
		Resources: resources,
	}

	data, err := yaml.Marshal(kustomize)
	if err != nil {
		return fmt.Errorf("failed to marshal kustomize resource: %w", err)
	}
	filename := filepath.Join(dir, "kustomization.yaml")
	if err = os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

func (c *ArtifactBuilder2) writeResourceWithName(obj runtime.Object, filename string) error {
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
	if err := e.Encode(obj, f); err != nil {
		return err
	}
	return nil
}

func (c *ArtifactBuilder2) writeResource(obj runtime.Object, dir string) error {
	name := obj.GetObjectKind().GroupVersionKind().Kind
	filename := filepath.Join(dir, fmt.Sprintf("%s.%s", name, "yaml"))
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
	if err := e.Encode(obj, f); err != nil {
		return err
	}
	return nil
}

func (c *ArtifactBuilder2) copyArtifacts(a artifact.Artifact, subDir, destDir string, repoLocations map[string]string) error {
	tmp, ok := repoLocations[a.ProfileRepoKey]
	if !ok {
		return fmt.Errorf("could not find repo clone for %q", a.ProfileRepoKey)
	}

	srcDir := filepath.Join(tmp, a.ProfilePath, subDir)
	if err := copy.Copy(srcDir, destDir); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}
	return nil
}

func (c *ArtifactBuilder2) makeHelmReleaseObjects(artifact profilesv1.Artifact, installation profilesv1.ProfileInstallation, definitionName string) (*helmv2.HelmRelease, *corev1.ConfigMap) {
	var helmChartSpec helmv2.HelmChartTemplateSpec
	if artifact.Chart.Path != "" {
		helmChartSpec = c.makeGitChartSpec(path.Join(installation.Spec.Source.Path, artifact.Chart.Path))
	} else if artifact.Chart != nil {
		helmChartSpec = c.makeHelmChartSpec(artifact.Chart.Name, artifact.Chart.Version, installation)
	}
	var (
		cfgMap *corev1.ConfigMap
		values []helmv2.ValuesReference
	)
	if artifact.Chart.DefaultValues != "" {
		cfgMap = c.makeDefaultValuesCfgMap(artifact.Name, artifact.Chart.DefaultValues, installation)
		// the default values always need to be at index 0
		values = append(values, helmv2.ValuesReference{
			Kind:      "ConfigMap",
			Name:      cfgMap.Name,
			ValuesKey: defaultValuesKey,
		})
	}
	if installation.Spec.ConfigMap != "" {
		artifactNameParts := strings.Split(artifact.Name, "/")
		values = append(values, helmv2.ValuesReference{
			Kind:      "ConfigMap",
			Name:      installation.Spec.ConfigMap,
			ValuesKey: artifactNameParts[len(artifactNameParts)-1],
		})
	}
	helmRelease := &helmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.makeArtifactName(artifact.Name, installation.Name, definitionName),
			Namespace: installation.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       helmv2.HelmReleaseKind,
			APIVersion: helmv2.GroupVersion.String(),
		},
		Spec: helmv2.HelmReleaseSpec{
			Chart: helmv2.HelmChartTemplate{
				Spec: helmChartSpec,
			},
			ValuesFrom: values,
		},
	}
	return helmRelease, cfgMap
}

func (c *ArtifactBuilder2) makeHelmRepository(url string, name string, installation profilesv1.ProfileInstallation) *sourcev1.HelmRepository {
	return &sourcev1.HelmRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.makeHelmRepoName(name, installation),
			Namespace: installation.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1.HelmRepositoryKind,
			APIVersion: sourcev1.GroupVersion.String(),
		},
		Spec: sourcev1.HelmRepositorySpec{
			URL: url,
		},
	}
}

func (c *ArtifactBuilder2) makeHelmRepoName(name string, installation profilesv1.ProfileInstallation) string {
	repoParts := strings.Split(installation.Spec.Source.URL, "/")
	repoName := repoParts[len(repoParts)-1]
	return c.join(installation.Name, repoName, name)
}

func (c *ArtifactBuilder2) makeGitChartSpec(path string) helmv2.HelmChartTemplateSpec {
	return helmv2.HelmChartTemplateSpec{
		Chart: path,
		SourceRef: helmv2.CrossNamespaceObjectReference{
			Kind:      sourcev1.GitRepositoryKind,
			Name:      c.GitRepositoryName,
			Namespace: c.GitRepositoryNamespace,
		},
	}
}

func (c *ArtifactBuilder2) makeHelmChartSpec(chart string, version string, installation profilesv1.ProfileInstallation) helmv2.HelmChartTemplateSpec {
	return helmv2.HelmChartTemplateSpec{
		Chart: chart,
		SourceRef: helmv2.CrossNamespaceObjectReference{
			Kind:      sourcev1.HelmRepositoryKind,
			Name:      c.makeHelmRepoName(chart, installation),
			Namespace: installation.ObjectMeta.Namespace,
		},
		Version: version,
	}
}

func (c *ArtifactBuilder2) makeDefaultValuesCfgMap(name, data string, installation profilesv1.ProfileInstallation) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.makeCfgMapName(name, installation),
			Namespace: installation.ObjectMeta.Namespace,
		},
		Data: map[string]string{
			defaultValuesKey: data,
		},
	}
}

func (c *ArtifactBuilder2) makeKustomizeHelmReleaseWrapper(artifact artifact.Artifact, installation profilesv1.ProfileInstallation, definitionName, path string, dependencies []artifact.Artifact) *kustomizev1.Kustomization {
	wrapper := c.makeKustomization(artifact, path, installation, definitionName, dependencies)
	wrapper.Spec.HealthChecks = []meta.NamespacedObjectKindReference{
		{
			APIVersion: helmv2.GroupVersion.String(),
			Kind:       helmv2.HelmReleaseKind,
			Name:       c.makeArtifactName(artifact.Name, installation.Name, definitionName),
			Namespace:  installation.ObjectMeta.Namespace,
		},
	}
	return wrapper
}

// makeKustomization creates a Kustomize object.
func (c *ArtifactBuilder2) makeKustomization(
	artifact artifact.Artifact,
	repoPath string,
	installation profilesv1.ProfileInstallation,
	definitionName string,
	dependencies []artifact.Artifact) *kustomizev1.Kustomization {
	var dependsOn []dependency.CrossNamespaceDependencyReference
	for _, dep := range dependencies {
		dependsOn = append(dependsOn, dependency.CrossNamespaceDependencyReference{
			Name:      c.makeArtifactName(dep.Name, installation.Name, definitionName),
			Namespace: installation.Namespace,
		})
	}
	return &kustomizev1.Kustomization{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.makeArtifactName(artifact.Name, installation.Name, definitionName),
			Namespace: installation.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		Spec: kustomizev1.KustomizationSpec{
			Path:            repoPath,
			Interval:        metav1.Duration{Duration: time.Minute * 5},
			Prune:           true,
			TargetNamespace: installation.ObjectMeta.Namespace,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      sourcev1.GitRepositoryKind,
				Name:      c.GitRepositoryName,
				Namespace: c.GitRepositoryNamespace,
			},
			DependsOn: dependsOn,
		},
	}
}

func (c *ArtifactBuilder2) makeCfgMapName(name string, installation profilesv1.ProfileInstallation) string {
	if strings.Contains(name, "/") {
		name = filepath.Base(name)
	}
	return c.join(installation.Name, name, "defaultvalues")
}

// join creates a joined string of name using - as a join character.
func (c *ArtifactBuilder2) join(s ...string) string {
	return strings.Join(s, "-")
}

// makeArtifactName creates a name for an artifact.
func (c *ArtifactBuilder2) makeArtifactName(name string, installationName, definitionName string) string {
	// if this is a nested artifact, it's name contains a /
	if strings.Contains(name, "/") {
		name = filepath.Base(name)
	}
	return c.join(installationName, definitionName, name)
}
func containsArtifact(list []artifact.Artifact, name string) (artifact.Artifact, bool) {
	for _, a := range list {
		if a.Name == name || a.ParentProfileArtifactName == name {
			return a, true
		}
	}
	return artifact.Artifact{}, false
}
