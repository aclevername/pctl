package install_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/otiai10/copy"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/api/types"

	fakegit "github.com/weaveworks/pctl/pkg/git/fakes"
	"github.com/weaveworks/pctl/pkg/install"
	"github.com/weaveworks/pctl/pkg/install/artifact"
)

var _ = Describe("Installer", func() {
	var (
		pSub                   profilesv1.ProfileInstallation
		fakeGitClient          *fakegit.FakeGit
		rootDir                string
		gitRepositoryNamespace string
		gitRepositoryName      string
	)

	BeforeEach(func() {
		pSub = profilesv1.ProfileInstallation{
			TypeMeta: profileTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name:      installationName,
				Namespace: namespace,
			},
			Spec: profilesv1.ProfileInstallationSpec{
				Source: &profilesv1.Source{
					URL:    profileURL,
					Branch: branch,
					Path:   profileName1,
				},
			},
		}
		fakeGitClient = &fakegit.FakeGit{}
		fakeGitClient.CloneStub = func(url string, branch string, dir string) error {
			from := filepath.Join("testdata", "simple_with_nested")
			err := copy.Copy(from, filepath.Join(dir))
			Expect(err).NotTo(HaveOccurred())
			return nil
		}
		var err error
		rootDir, err = ioutil.TempDir("", "test_make_artifacts")
		Expect(err).NotTo(HaveOccurred())
		gitRepositoryName = "git-repo-name"
		gitRepositoryNamespace = "git-repo-namespace"
		artifacts := []artifact.Artifact{
			{
				Objects: []artifact.Object{{
					Object: &helmv2.HelmRelease{
						TypeMeta: metav1.TypeMeta{
							Kind:       "HelmRelease",
							APIVersion: "helm.toolkit.fluxcd.io/v2beta1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-profile-weaveworks-nginx-dokuwiki",
							Namespace: "default",
						},
						Spec: helmv2.HelmReleaseSpec{
							Chart: helmv2.HelmChartTemplate{
								Spec: helmv2.HelmChartTemplateSpec{
									Chart:   "dokuwiki",
									Version: "11.1.6",
									SourceRef: helmv2.CrossNamespaceObjectReference{
										Kind:      "HelmRepository",
										Name:      "test-profile-profiles-examples-dokuwiki",
										Namespace: "default",
									},
								},
							},
							ValuesFrom: []helmv2.ValuesReference{
								{
									Name:     "nginx-values",
									Kind:     "Secret",
									Optional: true,
								},
							},
						},
					}, Path: "helm-chart"}, {
					Object: &kustomizev1.Kustomization{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-profile-weaveworks-nginx-dokuwiki",
							Namespace: "default",
						},
						TypeMeta: metav1.TypeMeta{
							Kind:       kustomizev1.KustomizationKind,
							APIVersion: kustomizev1.GroupVersion.String(),
						},
						Spec: kustomizev1.KustomizationSpec{
							Path:            "root-dir/artifacts/dokuwiki/helm-chart",
							Prune:           true,
							Interval:        metav1.Duration{Duration: time.Minute * 5},
							TargetNamespace: "default",
							SourceRef: kustomizev1.CrossNamespaceSourceReference{
								Kind:      sourcev1.GitRepositoryKind,
								Name:      gitRepositoryName,
								Namespace: gitRepositoryNamespace,
							},
							HealthChecks: []meta.NamespacedObjectKindReference{
								{
									APIVersion: helmv2.GroupVersion.String(),
									Kind:       helmv2.HelmReleaseKind,
									Name:       "test-profile-weaveworks-nginx-dokuwiki",
									Namespace:  "default",
								},
							},
						},
					}, Name: "kustomize-flux"}},
				Name:         "test-artifact-1",
				RepoURL:      "https://repo-url.com",
				PathsToCopy:  []string{"nginx/chart"},
				SparseFolder: "bitnami-nginx",
				Branch:       "",
				Kustomize: artifact.Kustomize{
					LocalResourceLimiter: &types.Kustomization{
						Resources: []string{"HelmRelease.yaml"},
					},
					ObjectWrapper: &types.Kustomization{
						Resources: []string{"kustomize-flux.yaml"},
					},
				},
				SubFolder: "helm-chart",
			},
			{
				Objects: []artifact.Object{{Object: &kustomizev1.Kustomization{
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
				}, Name: "kustomize-flux"}},
				Name:         "test-artifact-2",
				RepoURL:      "https://repo-url.com",
				PathsToCopy:  []string{"nginx/deployment"},
				SparseFolder: "weaveworks-nginx",
				Branch:       "",
				Kustomize: artifact.Kustomize{
					ObjectWrapper: &types.Kustomization{
						Resources: []string{"kustomize-flux.yaml"},
					},
				},
			},
		}
		install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
			return artifacts, nil
		})
	})

	AfterEach(func() {
		Expect(os.RemoveAll(rootDir)).To(Succeed())
	})

	Context("Install", func() {
		It("generates the artifacts", func() {
			installer := install.NewInstaller(install.Config{
				GitClient:        fakeGitClient,
				RootDir:          rootDir,
				GitRepoNamespace: gitRepoNamespace,
				GitRepoName:      gitRepoName,
			})
			err := installer.Install(pSub)
			Expect(err).NotTo(HaveOccurred())
			files := make(map[string]string)
			err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files[fmt.Sprintf("%s/%s", filepath.Base(filepath.Dir(path)), filepath.Base(path))] = path
				}
				return nil
			})
			Expect(err).NotTo(HaveOccurred())
			consistsOf := []string{
				filepath.Join(rootDir, "profile-installation.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-1", "helm-chart", "HelmRelease.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-1", "helm-chart", "kustomization.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-1", "helm-chart", "nginx", "chart.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-1", "kustomization.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-1", "kustomize-flux.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-2", "kustomization.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-2", "kustomize-flux.yaml"),
				filepath.Join(rootDir, "artifacts", "test-artifact-2", "nginx", "deployment", "deployment.yaml"),
			}
			Expect(files).To(ConsistOf(consistsOf))

			By("inspecting the profile-installation.yaml", func() {
				parent := filepath.Base(rootDir)
				content, err := ioutil.ReadFile(files[filepath.Join(parent, "profile-installation.yaml")])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(`apiVersion: weave.works/v1alpha1
kind: ProfileInstallation
metadata:
  creationTimestamp: null
  name: mySub
  namespace: default
spec:
  source:
    branch: main
    path: weaveworks-nginx
    url: https://github.com/org/repo-name
status: {}
`))
			})

			By("generating the path based kustomization artifact", func() {
				content, err := ioutil.ReadFile(files["test-artifact-1/kustomize-flux.yaml"])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(`apiVersion: kustomize.toolkit.fluxcd.io/v1beta1
kind: Kustomization
metadata:
  creationTimestamp: null
  name: test-profile-weaveworks-nginx-dokuwiki
  namespace: default
spec:
  healthChecks:
  - apiVersion: helm.toolkit.fluxcd.io/v2beta1
    kind: HelmRelease
    name: test-profile-weaveworks-nginx-dokuwiki
    namespace: default
  interval: 5m0s
  path: root-dir/artifacts/dokuwiki/helm-chart
  prune: true
  sourceRef:
    kind: GitRepository
    name: git-repo-name
    namespace: git-repo-namespace
  targetNamespace: default
status: {}
`))
			})

			By("generating a remote helm chart", func() {
				content, err := ioutil.ReadFile(files["helm-chart/HelmRelease.yaml"])
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(`apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: test-profile-weaveworks-nginx-dokuwiki
  namespace: default
spec:
  chart:
    spec:
      chart: dokuwiki
      sourceRef:
        kind: HelmRepository
        name: test-profile-profiles-examples-dokuwiki
        namespace: default
      version: 11.1.6
  interval: 0s
  valuesFrom:
  - kind: Secret
    name: nginx-values
    optional: true
status: {}
`))
			})
		})
		When("the profiles maker fails", func() {
			It("returns an error", func() {
				install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
					return nil, errors.New("nope")
				})
				installer := install.NewInstaller(install.Config{
					GitClient:        fakeGitClient,
					RootDir:          rootDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err := installer.Install(pSub)
				Expect(err).To(MatchError("failed to build artifact: nope"))
			})
		})
		When("fetching the local repository objects fails", func() {
			It("returns an error", func() {
				fakeGitClient.CloneReturns(errors.New("nope"))
				installer := install.NewInstaller(install.Config{
					GitClient:        fakeGitClient,
					RootDir:          rootDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err := installer.Install(pSub)
				Expect(err).To(MatchError("failed to get package local artifacts: failed to sparse clone folder with url: https://repo-url.com; branch: ; with error: nope"))
			})
		})
		When("there is a single profile repository", func() {
			It("creates files for all artifacts", func() {
				artifacts := []artifact.Artifact{
					{
						Objects: []artifact.Object{{Object: &helmv2.HelmRelease{
							TypeMeta: metav1.TypeMeta{
								Kind:       "HelmRelease",
								APIVersion: "helm.toolkit.fluxcd.io/v2beta1",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "test-profile-weaveworks-nginx-dokuwiki",
								Namespace: "default",
							},
							Spec: helmv2.HelmReleaseSpec{
								Chart: helmv2.HelmChartTemplate{
									Spec: helmv2.HelmChartTemplateSpec{
										Chart:   "dokuwiki",
										Version: "11.1.6",
										SourceRef: helmv2.CrossNamespaceObjectReference{
											Kind:      "HelmRepository",
											Name:      "test-profile-profiles-examples-dokuwiki",
											Namespace: "default",
										},
									},
								},
								Values: &apiextensionsv1.JSON{
									Raw: []byte(`{"replicaCount": 3,"service":{"port":8081}}`),
								},
								ValuesFrom: []helmv2.ValuesReference{
									{
										Name:     "nginx-values",
										Kind:     "Secret",
										Optional: true,
									},
								},
							},
						}, Path: "helm-chart"}},
						Name:         "test-artifact-1",
						RepoURL:      "https://repo-url.com",
						PathsToCopy:  []string{"nginx/chart"},
						SparseFolder: "bitnami-nginx",
						Branch:       "",
					},
				}
				install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
					return artifacts, nil
				})
				pSub := profilesv1.ProfileInstallation{
					TypeMeta: profileTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      installationName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileInstallationSpec{
						Source: &profilesv1.Source{
							URL:    "https://github.com/weaveworks/nginx-profile",
							Branch: "main",
						},
					},
				}
				tempDir, err := ioutil.TempDir("", "catalog-install")
				Expect(err).NotTo(HaveOccurred())
				installer := install.NewInstaller(install.Config{
					ProfileName:      "generate-test",
					GitClient:        fakeGitClient,
					RootDir:          filepath.Join(tempDir, "generate-test"),
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err = installer.Install(pSub)
				Expect(err).NotTo(HaveOccurred())

				var files []string
				err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						files = append(files, path)
					}
					return nil
				})
				Expect(err).NotTo(HaveOccurred())

				profileFile := filepath.Join(tempDir, "generate-test", "profile-installation.yaml")
				artifactHelmRelease := filepath.Join(tempDir, "generate-test", "artifacts", "test-artifact-1", "helm-chart", "HelmRelease.yaml")

				Expect(hasCorrectFilePerms(profileFile)).To(BeTrue())
				Expect(hasCorrectFilePerms(artifactHelmRelease)).To(BeTrue())

				content, err := ioutil.ReadFile(profileFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(Equal(`apiVersion: weave.works/v1alpha1
kind: ProfileInstallation
metadata:
  creationTimestamp: null
  name: mySub
  namespace: default
spec:
  source:
    branch: main
    url: https://github.com/weaveworks/nginx-profile
status: {}
`))

				content, err = ioutil.ReadFile(artifactHelmRelease)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(Equal(`apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: test-profile-weaveworks-nginx-dokuwiki
  namespace: default
spec:
  chart:
    spec:
      chart: dokuwiki
      sourceRef:
        kind: HelmRepository
        name: test-profile-profiles-examples-dokuwiki
        namespace: default
      version: 11.1.6
  interval: 0s
  values:
    replicaCount: 3
    service:
      port: 8081
  valuesFrom:
  - kind: Secret
    name: nginx-values
    optional: true
status: {}
`))
			})
		})

		When("there are multiple artifacts for the same repository", func() {
			var artifacts []artifact.Artifact
			BeforeEach(func() {
				artifacts = []artifact.Artifact{
					{
						Objects: []artifact.Object{
							{Object: &kustomizev1.Kustomization{
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
							}, Path: "helm-chart"},
						},
						Name:         "test-artifact-1",
						RepoURL:      "https://repo-url.com",
						PathsToCopy:  []string{"nginx/deployment"},
						SparseFolder: "weaveworks-nginx-1",
						Branch:       "main",
					},
					{
						Objects: []artifact.Object{
							{Object: &kustomizev1.Kustomization{
								TypeMeta: metav1.TypeMeta{
									Kind:       "Kustomization",
									APIVersion: "kustomize.toolkit.fluxcd.io/v1beta1",
								},
								ObjectMeta: metav1.ObjectMeta{
									Name:      "test-profile-weaveworks-nginx-kustomize-2",
									Namespace: "default",
								},
								Spec: kustomizev1.KustomizationSpec{
									Path: "root-dir/artifacts/kustomize/nginx2/deployment",
									SourceRef: kustomizev1.CrossNamespaceSourceReference{
										Kind:      "GitRepository",
										Namespace: gitRepositoryNamespace,
										Name:      gitRepositoryName,
									},
									Interval:        metav1.Duration{Duration: 300000000000},
									Prune:           true,
									TargetNamespace: "default",
								},
							}, Path: "helm-chart"},
						},
						Name:         "test-artifact-2",
						RepoURL:      "https://repo-url.com",
						PathsToCopy:  []string{"nginx2/deployment"},
						SparseFolder: "weaveworks-nginx-2",
						Branch:       "main",
					},
				}
				pSub = profilesv1.ProfileInstallation{
					TypeMeta: profileTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      installationName,
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileInstallationSpec{
						Source: &profilesv1.Source{
							URL:    "https://github.com/weaveworks/nginx-profile",
							Branch: "main",
						},
					},
				}
				fakeGitClient.CloneStub = func(url string, branch string, dir string) error {
					from := filepath.Join("testdata", "clone_cache")
					err := copy.Copy(from, filepath.Join(dir))
					Expect(err).NotTo(HaveOccurred())
					return nil
				}
			})

			It("shouldn't clone as many times as there are artifacts", func() {
				fakeGitClient.CloneStub = func(url string, branch string, dir string) error {
					from := filepath.Join("testdata", "clone_cache")
					err := copy.Copy(from, filepath.Join(dir))
					Expect(err).NotTo(HaveOccurred())
					return nil
				}
				install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
					return artifacts, nil
				})
				installer := install.NewInstaller(install.Config{
					GitClient:        fakeGitClient,
					RootDir:          rootDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err := installer.Install(pSub)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeGitClient.CloneCallCount()).To(Equal(1), "cache did not work, clone has been called more than once.")
			})

			It("should still clone different urls", func() {
				fakeGitClient.CloneStub = func(url string, branch string, dir string) error {
					from := filepath.Join("testdata", "clone_cache")
					err := copy.Copy(from, filepath.Join(dir))
					Expect(err).NotTo(HaveOccurred())
					return nil
				}
				artifacts[0].RepoURL = "https://github.com/weaveworks/profiles-examples"
				install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
					return artifacts, nil
				})
				installer := install.NewInstaller(install.Config{
					GitClient:        fakeGitClient,
					RootDir:          rootDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err := installer.Install(pSub)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeGitClient.CloneCallCount()).To(Equal(2))
			})

			It("should still clone same urls but with different branch", func() {
				fakeGitClient.CloneStub = func(url string, branch string, dir string) error {
					from := filepath.Join("testdata", "clone_cache")
					err := copy.Copy(from, filepath.Join(dir))
					Expect(err).NotTo(HaveOccurred())
					return nil
				}
				artifacts[0].Branch = "different"
				install.SetProfileMakeArtifacts(func(i *install.Installer, installation profilesv1.ProfileInstallation) ([]artifact.Artifact, error) {
					return artifacts, nil
				})
				installer := install.NewInstaller(install.Config{
					GitClient:        fakeGitClient,
					RootDir:          rootDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				})
				err := installer.Install(pSub)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeGitClient.CloneCallCount()).To(Equal(2))
			})
		})
	})
})

func hasCorrectFilePerms(file string) bool {
	info, err := os.Stat(file)
	Expect(err).NotTo(HaveOccurred())
	return strconv.FormatUint(uint64(info.Mode().Perm()), 8) == strconv.FormatInt(0644, 8)
}
