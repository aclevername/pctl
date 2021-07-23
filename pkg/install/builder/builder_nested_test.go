package builder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/runtime/dependency"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/pctl/pkg/install/builder"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Builder", func() {
	Context("when the artifact is a nested artifact", func() {
		BeforeEach(func() {
			kustomizeFilesDir := filepath.Join(gitDir, profilePath, "files")
			Expect(os.MkdirAll(kustomizeFilesDir, 0755)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(kustomizeFilesDir, "file1"), []byte("foo"), 0755)).To(Succeed())
			artifacts = []builder.ArtifactWrapper{
				{
					Artifact: profilesv1.Artifact{
						Name: artifactName,
						Kustomize: &profilesv1.Kustomize{
							Path: "files/",
						},
					},
					PathToProfileClone:        filepath.Join(gitDir, profilePath),
					ProfileName:               profileName,
					NestedProfileArtifactName: "nested-profile",
				},
			}
		})

		It("places the artifactr in a subdirectory", func() {
			err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
			Expect(err).NotTo(HaveOccurred())

			var files []string
			err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, strings.TrimPrefix(path, rootDir+"/"))
				}
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(files).To(ConsistOf(
				"artifacts/nested-profile/1/kustomization.yaml",
				"artifacts/nested-profile/1/kustomize-flux.yaml",
				"artifacts/nested-profile/1/files/file1",
				"profile-installation.yaml",
			))
		})
	})

	Context("dependsOn", func() {
		var (
			artifactName2     = "2"
			artifactName3     = "3"
			nestedProfileName = "nested-profile"
		)
		BeforeEach(func() {
			kustomizeFilesDir := filepath.Join(gitDir, profilePath, "files")
			Expect(os.MkdirAll(kustomizeFilesDir, 0755)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(kustomizeFilesDir, "file1"), []byte("foo"), 0755)).To(Succeed())
			artifacts = []builder.ArtifactWrapper{
				{
					Artifact: profilesv1.Artifact{
						Name: artifactName,
						Kustomize: &profilesv1.Kustomize{
							Path: "files/",
						},
						DependsOn: []profilesv1.DependsOn{
							{
								Name: artifactName2,
							},
							{
								Name: nestedProfileName,
							},
						},
					},
					PathToProfileClone: filepath.Join(gitDir, profilePath),
					ProfileName:        profileName,
				},
				{
					Artifact: profilesv1.Artifact{
						Name: artifactName2,
						Kustomize: &profilesv1.Kustomize{
							Path: "files/",
						},
					},
					PathToProfileClone: filepath.Join(gitDir, profilePath),
					ProfileName:        profileName,
				},
				{
					Artifact: profilesv1.Artifact{
						Name: artifactName3,
						Kustomize: &profilesv1.Kustomize{
							Path: "files/",
						},
					},
					NestedProfileArtifactName: nestedProfileName,
					PathToProfileClone:        filepath.Join(gitDir, profilePath),
					ProfileName:               profileName,
				},
			}
		})

		It("sets the depends on field in the kustomization", func() {
			err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
			Expect(err).NotTo(HaveOccurred())

			var files []string
			err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, strings.TrimPrefix(path, rootDir+"/"))
				}
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(files).To(ConsistOf(
				"artifacts/1/kustomization.yaml",
				"artifacts/1/kustomize-flux.yaml",
				"artifacts/1/files/file1",
				"artifacts/2/kustomization.yaml",
				"artifacts/2/kustomize-flux.yaml",
				"artifacts/2/files/file1",
				"artifacts/nested-profile/3/kustomization.yaml",
				"artifacts/nested-profile/3/kustomize-flux.yaml",
				"artifacts/nested-profile/3/files/file1",
				"profile-installation.yaml",
			))

			kustomize := kustomizev1.Kustomization{}
			decodeFile(filepath.Join(rootDir, "artifacts/1/kustomize-flux.yaml"), &kustomize)
			Expect(kustomize).To(Equal(kustomizev1.Kustomization{
				TypeMeta: kustomizeTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s-%s", installationName, profileName, artifactName),
					Namespace: namespace,
				},
				Spec: kustomizev1.KustomizationSpec{
					Path: filepath.Join(rootDir, "artifacts/1/files/"),
					SourceRef: kustomizev1.CrossNamespaceSourceReference{
						Kind:      "GitRepository",
						Namespace: gitRepoNamespace,
						Name:      gitRepoName,
					},
					Interval:        metav1.Duration{Duration: 300000000000},
					Prune:           true,
					TargetNamespace: namespace,
					DependsOn: []dependency.CrossNamespaceDependencyReference{
						{
							Name:      fmt.Sprintf("%s-%s-%s", installationName, profileName, artifactName2),
							Namespace: namespace,
						},
						{
							Name:      fmt.Sprintf("%s-%s-%s", installationName, profileName, artifactName3),
							Namespace: namespace,
						},
					},
				},
			}))
		})
	})

	Context("invalid artifacts", func() {
		When("no type is set", func() {
			It("returns an error", func() {
				artifacts[0].Artifact = profilesv1.Artifact{}
				err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
				Expect(err).To(MatchError(ContainSubstring("no artifact type set")))
			})
		})

		When("remote chart and path is configured", func() {
			It("returns an error", func() {
				artifacts[0].Artifact = profilesv1.Artifact{
					Chart: &profilesv1.Chart{
						Path: "foo",
						URL:  "bar",
					},
				}
				err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
				Expect(err).To(MatchError(ContainSubstring("expected exactly one, got both: chart.path, chart.url")))
			})
		})

		When("chart and kustomize is configured", func() {
			It("returns an error", func() {
				artifacts[0].Artifact = profilesv1.Artifact{
					Chart: &profilesv1.Chart{
						Path: "foo",
					},
					Kustomize: &profilesv1.Kustomize{
						Path: "bar",
					},
				}
				err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
				Expect(err).To(MatchError(ContainSubstring("expected exactly one, got both: chart, kustomize")))
			})
		})

		When("chart and profile is configured", func() {
			It("returns an error", func() {
				artifacts[0].Artifact = profilesv1.Artifact{
					Profile: &profilesv1.Profile{},
					Chart: &profilesv1.Chart{
						Path: "foo",
					},
				}
				err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
				Expect(err).To(MatchError(ContainSubstring("expected exactly one, got both: chart, profile")))
			})
		})

		When("kustomize and profile is configured", func() {
			It("returns an error", func() {
				artifacts[0].Artifact = profilesv1.Artifact{
					Profile: &profilesv1.Profile{},
					Kustomize: &profilesv1.Kustomize{
						Path: "foo",
					},
				}
				err := artifactBuilder.Write(installation, artifacts, repoLocationMap)
				Expect(err).To(MatchError(ContainSubstring("expected exactly one, got both: kustomize, profile")))
			})
		})
	})
})
