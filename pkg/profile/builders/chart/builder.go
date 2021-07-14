package chart

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/weaveworks/pctl/pkg/profile/artifact"
)

const defaultValuesKey = "default-values.yaml"

// Config defines some common configuration values for builders.
type Config struct {
	GitRepositoryName      string
	GitRepositoryNamespace string
	RootDir                string
}

// Builder will build helm chart resources.
type Builder struct {
	Config
}

// Build a single artifact from a profile artifact and installation.
func (c *Builder) Build(att profilesv1.Artifact, installation profilesv1.ProfileInstallation, definition profilesv1.ProfileDefinition) ([]artifact.Artifact, error) {
	if err := validateArtifact(att); err != nil {
		return nil, fmt.Errorf("validation failed for artifact %s: %w", att.Name, err)
	}
	// but we want the kustomization outside this so return multiple artifacts.
	a := artifact.Artifact{Name: att.Name, SubFolder: "helm-chart"}
	//var deps []profilesv1.Artifact
	//for _, dep := range att.DependsOn {
	//	d, ok := containsArtifact(dep.Name, definition.Spec.Artifacts)
	//	if !ok {
	//		return nil, fmt.Errorf("%s's depending artifact %s not found in the list of artifacts", a.Name, dep.Name)
	//	}
	//
	//	deps = append(deps, d)
	//}

	helmRelease, cfgMap := c.makeHelmReleaseObjects(att, installation, definition.Name)
	if cfgMap != nil {
		a.Objects = append(a.Objects, cfgMap)
	}
	a.Objects = append(a.Objects, helmRelease)
	if att.Chart.Path != "" {
		if c.GitRepositoryNamespace == "" && c.GitRepositoryName == "" {
			return nil, fmt.Errorf("in case of local resources, the flux gitrepository object's details must be provided")
		}
		helmRelease.Spec.Chart.Spec.Chart = filepath.Join(c.RootDir, "artifacts", "helm-chart", att.Name, att.Chart.Path)
		branch := installation.Spec.Source.Branch
		if installation.Spec.Source.Tag != "" {
			branch = installation.Spec.Source.Tag
		}
		a.RepoURL = installation.Spec.Source.URL
		a.SparseFolder = definition.Name
		a.Branch = branch
		a.PathsToCopy = append(a.PathsToCopy, att.Chart.Path)
		a.Kustomize = &types.Kustomization{
			Resources: []string{"HelmRelease.yaml"},
		}
	}
	if att.Chart.URL != "" {
		helmRepository := c.makeHelmRepository(att.Chart.URL, att.Chart.Name, installation)
		a.Objects = append(a.Objects, helmRepository)
	}
	a.HelmWrapper = &types.Kustomization{
		Resources: []string{"kustomize-flux.yaml"},
	}
	a.HelmWrapperKustomization = c.makeKustomizeWrapper(att, installation, definition.Name)
	return []artifact.Artifact{a}, nil
}

// validateArtifact validates that the artifact has valid chart properties.
func validateArtifact(in profilesv1.Artifact) error {
	if in.Profile != nil {
		return apis.ErrMultipleOneOf("chart", "profile")
	}
	if in.Kustomize != nil {
		return apis.ErrMultipleOneOf("chart", "kustomize")
	}
	if in.Chart.Path != "" && in.Chart.URL != "" {
		return apis.ErrMultipleOneOf("chart.path", "chart.url")
	}
	return nil
}

func (c *Builder) makeHelmReleaseObjects(artifact profilesv1.Artifact, installation profilesv1.ProfileInstallation, definitionName string) (*helmv2.HelmRelease, *corev1.ConfigMap) {
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
			Name:      makeArtifactName(artifact.Name, installation.Name, definitionName),
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

func (c *Builder) makeHelmRepository(url string, name string, installation profilesv1.ProfileInstallation) *sourcev1.HelmRepository {
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

func (c *Builder) makeHelmRepoName(name string, installation profilesv1.ProfileInstallation) string {
	repoParts := strings.Split(installation.Spec.Source.URL, "/")
	repoName := repoParts[len(repoParts)-1]
	return join(installation.Name, repoName, name)
}

// TODO: this will not create it with kustomize-flux -> because it uses the Kind not the name.
func (c *Builder) makeKustomizeWrapper(artifact profilesv1.Artifact, installation profilesv1.ProfileInstallation, definitionName string) *kustomizev1.Kustomization {
	path := filepath.Join(c.RootDir, "artifacts", artifact.Name, "helm-chart")
	return &kustomizev1.Kustomization{
		ObjectMeta: metav1.ObjectMeta{
			Name:      makeArtifactName(artifact.Name+"-wrapper", installation.Name, definitionName),
			Namespace: installation.ObjectMeta.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		Spec: kustomizev1.KustomizationSpec{
			Path:            path,
			Interval:        metav1.Duration{Duration: time.Minute * 5},
			Prune:           true,
			TargetNamespace: installation.ObjectMeta.Namespace,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      sourcev1.GitRepositoryKind,
				Name:      c.GitRepositoryName,
				Namespace: c.GitRepositoryNamespace,
			},
			//DependsOn: dependsOn,
			HealthChecks: []meta.NamespacedObjectKindReference{
				{
					APIVersion: helmv2.GroupVersion.String(),
					Kind:       helmv2.HelmReleaseKind,
					Name:       makeArtifactName(artifact.Name, installation.Name, definitionName),
					Namespace:  installation.ObjectMeta.Namespace,
				},
			},
		},
	}
}

func makeArtifactName(name string, installationName, definitionName string) string {
	// if this is a nested artifact, it's name contains a /
	if strings.Contains(name, "/") {
		name = filepath.Base(name)
	}
	return join(installationName, definitionName, name)
}

func join(s ...string) string {
	return strings.Join(s, "-")
}

func (c *Builder) makeGitChartSpec(path string) helmv2.HelmChartTemplateSpec {
	return helmv2.HelmChartTemplateSpec{
		Chart: path,
		SourceRef: helmv2.CrossNamespaceObjectReference{
			Kind:      sourcev1.GitRepositoryKind,
			Name:      c.GitRepositoryName,
			Namespace: c.GitRepositoryNamespace,
		},
	}
}

func (c *Builder) makeHelmChartSpec(chart string, version string, installation profilesv1.ProfileInstallation) helmv2.HelmChartTemplateSpec {
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

func (c *Builder) makeDefaultValuesCfgMap(name, data string, installation profilesv1.ProfileInstallation) *corev1.ConfigMap {
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

func (c *Builder) makeCfgMapName(name string, installation profilesv1.ProfileInstallation) string {
	if strings.Contains(name, "/") {
		name = filepath.Base(name)
	}
	return join(installation.Name, name, "defaultvalues")
}

//// containsArtifact checks whether an artifact with a specific name exists in a list of artifacts.
//func containsArtifact(name string, stack []profilesv1.Artifact) (profilesv1.Artifact, bool) {
//	for _, a := range stack {
//		if a.Name == name {
//			return a, true
//		}
//	}
//	return profilesv1.Artifact{}, false
//}
