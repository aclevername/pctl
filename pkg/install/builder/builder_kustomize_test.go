package builder_test

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/runtime/dependency"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/weaveworks/pctl/pkg/install/artifact"
	"github.com/weaveworks/pctl/pkg/install/builder"
)

var _ = Describe("ArtifactBuilder", func() {
	var (
		profileName            string
		profileURL             string
		profilePath            string
		partifact              profilesv1.Artifact
		pSub                   profilesv1.ProfileInstallation
		pDef                   profilesv1.ProfileDefinition
		rootDir                string
		gitRepositoryName      string
		gitRepositoryNamespace string
		profileName1           = "weaveworks-nginx"
		namespace              = "default"
		profileSubAPIVersion   = "weave.works/v1alpha1"
		profileSubKind         = "ProfileInstallation"
	)

	var (
		profileTypeMeta = metav1.TypeMeta{
			Kind:       profileSubKind,
			APIVersion: profileSubAPIVersion,
		}
	)
	BeforeEach(func() {
		profileName = "test-profile"
		profileURL = "https://github.com/weaveworks/profiles-examples"
		profilePath = "weaveworks-nginx"
		pSub = profilesv1.ProfileInstallation{
			TypeMeta: profileTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      profileName,
				Namespace: namespace,
			},
			Spec: profilesv1.ProfileInstallationSpec{
				Source: &profilesv1.Source{
					URL:  profileURL,
					Tag:  "weaveworks-nginx/v0.0.1",
					Path: profilePath,
				},
			},
		}
		partifact = profilesv1.Artifact{
			Name: "kustomize",
			Kustomize: &profilesv1.Kustomize{
				Path: "nginx/deployment",
			},
		}
		pDef = profilesv1.ProfileDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: profileName1,
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "Profile",
				APIVersion: "packages.weave.works/profilesv1",
			},
			Spec: profilesv1.ProfileDefinitionSpec{
				ProfileDescription: profilesv1.ProfileDescription{
					Description: "foo",
				},
				Artifacts: []profilesv1.Artifact{partifact},
			},
		}
		rootDir = "root-dir"
		gitRepositoryName = "git-repository-name"
		gitRepositoryNamespace = "git-repository-namespace"
	})

	Context("Build", func() {
		It("creates an artifact from an install and a profile definition", func() {
			builder := &builder.ArtifactBuilder{
				Config: builder.Config{
					GitRepositoryName:      gitRepositoryName,
					GitRepositoryNamespace: gitRepositoryNamespace,
					RootDir:                rootDir,
				},
			}
			artifacts, err := builder.Build(partifact, pSub, pDef)
			Expect(err).NotTo(HaveOccurred())
			kustomization := &kustomizev1.Kustomization{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Kustomization",
					APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-profile-weaveworks-nginx-kustomize",
					Namespace: "default",
				},
				Spec: kustomizev1.KustomizationSpec{
					Path: "root-dir/artifacts/kustomize/nginx/deployment",
					SourceRef: kustomizev1.CrossNamespaceSourceReference{
						Kind:      "GitRepository",
						Namespace: gitRepositoryNamespace,
						Name:      gitRepositoryName,
					},
					Interval:        metav1.Duration{Duration: 300000000000},
					Prune:           true,
					TargetNamespace: "default",
				},
			}
			expected := artifact.Artifact{
				Objects:      []artifact.Object{{Object: kustomization, Name: "kustomize-flux"}},
				Name:         "kustomize",
				RepoURL:      "https://github.com/weaveworks/profiles-examples",
				PathsToCopy:  []string{"nginx/deployment"},
				SparseFolder: "weaveworks-nginx",
				Branch:       "weaveworks-nginx/v0.0.1",
				Kustomize: artifact.Kustomize{
					ObjectWrapper: &types.Kustomization{
						Resources: []string{"kustomize-flux.yaml"},
					},
				},
			}
			Expect(artifacts[0]).To(Equal(expected))
		})
		When("branch is defined instead of tag", func() {
			It("will use the branch definition", func() {
				pSub = profilesv1.ProfileInstallation{
					TypeMeta: profileTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      profileName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileInstallationSpec{
						Source: &profilesv1.Source{
							URL:    profileURL,
							Branch: "custom-branch",
							Path:   profilePath,
						},
					},
				}
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						GitRepositoryName:      gitRepositoryName,
						GitRepositoryNamespace: gitRepositoryNamespace,
						RootDir:                rootDir,
					},
				}
				artifacts, err := builder.Build(partifact, pSub, pDef)
				Expect(err).NotTo(HaveOccurred())
				kustomization := &kustomizev1.Kustomization{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Kustomization",
						APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-profile-weaveworks-nginx-kustomize",
						Namespace: "default",
					},
					Spec: kustomizev1.KustomizationSpec{
						Path: "root-dir/artifacts/kustomize/nginx/deployment",
						SourceRef: kustomizev1.CrossNamespaceSourceReference{
							Kind:      "GitRepository",
							Namespace: gitRepositoryNamespace,
							Name:      gitRepositoryName,
						},
						Interval:        metav1.Duration{Duration: 300000000000},
						Prune:           true,
						TargetNamespace: "default",
					},
				}
				expected := artifact.Artifact{
					Objects:      []artifact.Object{{Object: kustomization, Name: "kustomize-flux"}},
					Name:         "kustomize",
					RepoURL:      "https://github.com/weaveworks/profiles-examples",
					PathsToCopy:  []string{"nginx/deployment"},
					SparseFolder: "weaveworks-nginx",
					Branch:       "custom-branch",
					Kustomize: artifact.Kustomize{
						ObjectWrapper: &types.Kustomization{
							Resources: []string{"kustomize-flux.yaml"},
						},
					},
				}
				Expect(artifacts).To(ConsistOf(expected))
			})
		})
		When("git-repository-name and git-repository-namespace aren't defined", func() {
			It("returns an error", func() {
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						RootDir: rootDir,
					},
				}
				partifact = profilesv1.Artifact{
					Name: "local-partifact",
					Chart: &profilesv1.Chart{
						Path: "nginx/chart",
					},
				}
				pDef = profilesv1.ProfileDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name: profileName1,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Profile",
						APIVersion: "packages.weave.works/profilesv1",
					},
					Spec: profilesv1.ProfileDefinitionSpec{
						ProfileDescription: profilesv1.ProfileDescription{
							Description: "foo",
						},
						Artifacts: []profilesv1.Artifact{partifact},
					},
				}
				_, err := builder.Build(partifact, pSub, pDef)
				Expect(err).To(MatchError("in case of local resources, the flux gitrepository object's details must be provided"))
			})
		})
		When("profile and kustomize", func() {
			It("errors", func() {
				a := profilesv1.Artifact{
					Name: "test",
					Profile: &profilesv1.Profile{
						Source: &profilesv1.Source{
							URL:    "example.com",
							Branch: "branch",
						},
					},
					Kustomize: &profilesv1.Kustomize{
						Path: "https://not.empty",
					},
				}
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						RootDir:                rootDir,
						GitRepositoryNamespace: gitRepositoryNamespace,
						GitRepositoryName:      gitRepositoryName,
					},
				}
				_, err := builder.Build(a, pSub, pDef)
				Expect(err).To(MatchError(ContainSubstring("validation failed for artifact test: expected exactly one, got both: kustomize, profile")))
			})
		})
		When("chart and kustomize", func() {
			It("errors", func() {
				a := profilesv1.Artifact{
					Name: "test",
					Chart: &profilesv1.Chart{
						Name: "chart",
					},
					Kustomize: &profilesv1.Kustomize{
						Path: "https://not.empty",
					},
				}
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						RootDir:                rootDir,
						GitRepositoryNamespace: gitRepositoryNamespace,
						GitRepositoryName:      gitRepositoryName,
					},
				}
				_, err := builder.Build(a, pSub, pDef)
				Expect(err).To(MatchError(ContainSubstring("validation failed for artifact test: expected exactly one, got both: chart, kustomize")))
			})
		})
		When("depends on is defined for an artifact", func() {
			It("creates a kustomize object with DependsOn set correctly", func() {
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						GitRepositoryName:      gitRepositoryName,
						GitRepositoryNamespace: gitRepositoryNamespace,
						RootDir:                rootDir,
					},
				}
				partifact = profilesv1.Artifact{
					Name: "kustomize",
					Kustomize: &profilesv1.Kustomize{
						Path: "nginx/deployment",
					},
					DependsOn: []profilesv1.DependsOn{
						{
							Name: "depends-on",
						},
					},
				}
				partifacts := []profilesv1.Artifact{{
					Name: "depends-on",
					Kustomize: &profilesv1.Kustomize{
						Path: "nginx/deployment",
					},
				}, partifact}
				pDef = profilesv1.ProfileDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name: profileName1,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Profile",
						APIVersion: "packages.weave.works/profilesv1",
					},
					Spec: profilesv1.ProfileDefinitionSpec{
						ProfileDescription: profilesv1.ProfileDescription{
							Description: "foo",
						},
						Artifacts: partifacts,
					},
				}
				artifacts, err := builder.Build(partifact, pSub, pDef)
				Expect(err).NotTo(HaveOccurred())
				kustomization := &kustomizev1.Kustomization{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Kustomization",
						APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-profile-weaveworks-nginx-kustomize",
						Namespace: "default",
					},
					Spec: kustomizev1.KustomizationSpec{
						Path: "root-dir/artifacts/kustomize/nginx/deployment",
						SourceRef: kustomizev1.CrossNamespaceSourceReference{
							Kind:      "GitRepository",
							Namespace: gitRepositoryNamespace,
							Name:      gitRepositoryName,
						},
						Interval:        metav1.Duration{Duration: 300000000000},
						Prune:           true,
						TargetNamespace: "default",
						DependsOn: []dependency.CrossNamespaceDependencyReference{
							{
								Namespace: "default",
								Name:      "test-profile-weaveworks-nginx-depends-on",
							},
						},
					},
				}
				expected := artifact.Artifact{
					Objects:      []artifact.Object{{Object: kustomization, Name: "kustomize-flux"}},
					Name:         "kustomize",
					RepoURL:      "https://github.com/weaveworks/profiles-examples",
					PathsToCopy:  []string{"nginx/deployment"},
					SparseFolder: "weaveworks-nginx",
					Branch:       "weaveworks-nginx/v0.0.1",
					Kustomize: artifact.Kustomize{
						ObjectWrapper: &types.Kustomization{
							Resources: []string{"kustomize-flux.yaml"},
						},
					},
				}
				Expect(artifacts).To(ConsistOf(expected))
			})
		})
		When("depends on is defined for an artifact but the artifact is not in the list", func() {
			It("returns a sensible error", func() {
				builder := &builder.ArtifactBuilder{
					Config: builder.Config{
						GitRepositoryName:      gitRepositoryName,
						GitRepositoryNamespace: gitRepositoryNamespace,
						RootDir:                rootDir,
					},
				}
				partifact = profilesv1.Artifact{
					Name: "kustomize",
					Kustomize: &profilesv1.Kustomize{
						Path: "nginx/deployment",
					},
					DependsOn: []profilesv1.DependsOn{
						{
							Name: "depends-on",
						},
					},
				}
				pDef = profilesv1.ProfileDefinition{
					ObjectMeta: metav1.ObjectMeta{
						Name: profileName1,
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       "Profile",
						APIVersion: "packages.weave.works/profilesv1",
					},
					Spec: profilesv1.ProfileDefinitionSpec{
						ProfileDescription: profilesv1.ProfileDescription{
							Description: "foo",
						},
						Artifacts: []profilesv1.Artifact{partifact},
					},
				}
				_, err := builder.Build(partifact, pSub, pDef)
				Expect(err).To(MatchError("kustomize's depending artifact depends-on not found in the list of artifacts"))
			})
		})
	})
})
