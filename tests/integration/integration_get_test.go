package integration_test

import (
	"context"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("pctl get", func() {
	Context("get catalog profiles", func() {
		It("returns the matching profiles", func() {
			cmd := exec.Command(binaryPath, "get", "--catalog", "nginx")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())
			expected := "CATALOG/PROFILE               	VERSION	DESCRIPTION                     \n" +
				"nginx-catalog/weaveworks-nginx	v0.1.0 	This installs nginx.           \t\n" +
				"nginx-catalog/weaveworks-nginx	v0.1.1 	This installs nginx.           \t\n" +
				"nginx-catalog/bitnami-nginx   	v0.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/nginx           	v2.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/some-other-nginx	       	This installs some other nginx.\t\n\n"
			Expect(string(session)).To(ContainSubstring(expected))
		})

		It("returns all the catalog profiles", func() {
			cmd := exec.Command(binaryPath, "get", "--catalog")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())
			expected := "CATALOG/PROFILE               	VERSION	DESCRIPTION                     \n" +
				"nginx-catalog/weaveworks-nginx	v0.1.0 	This installs nginx.           \t\n" +
				"nginx-catalog/weaveworks-nginx	v0.1.1 	This installs nginx.           \t\n" +
				"nginx-catalog/bitnami-nginx   	v0.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/nginx           	v2.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/some-other-nginx	       	This installs some other nginx.\t\n\n"
			Expect(string(session)).To(ContainSubstring(expected))
		})

		When("-o is set to json", func() {
			It("returns the matching profiles in json", func() {
				cmd := exec.Command(binaryPath, "get", "-o", "json", "--catalog", "nginx")
				session, err := cmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(string(session)).To(ContainSubstring(`{
    "tag": "weaveworks-nginx/v0.1.0",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "weaveworks-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "weaveworks-nginx/v0.1.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "weaveworks-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "bitnami-nginx/v0.0.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "bitnami-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "v2.0.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/nginx-profile",
    "name": "nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "catalogSource": "nginx-catalog",
    "name": "some-other-nginx",
    "description": "This installs some other nginx."
  }`))
			})

			It("returns all catalog profiles in json", func() {
				cmd := exec.Command(binaryPath, "get", "-o", "json", "--catalog")
				session, err := cmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(string(session)).To(ContainSubstring(`{
    "tag": "weaveworks-nginx/v0.1.0",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "weaveworks-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "weaveworks-nginx/v0.1.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "weaveworks-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "bitnami-nginx/v0.0.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/profiles-examples",
    "name": "bitnami-nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "tag": "v2.0.1",
    "catalogSource": "nginx-catalog",
    "url": "https://github.com/weaveworks/nginx-profile",
    "name": "nginx",
    "description": "This installs nginx.",
    "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
    "prerequisites": [
      "Kubernetes 1.18+"
    ]
  },
  {
    "catalogSource": "nginx-catalog",
    "name": "some-other-nginx",
    "description": "This installs some other nginx."
  }`))
			})
		})

		It("doesnt not break json output", func() {
			cmd := exec.Command(binaryPath, "get", "-o", "json", "--catalog")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(session)).ToNot(ContainSubstring("PACKAGE CATALOG"))
		})

		When("kubeconfig is incorrectly set", func() {
			It("returns a useful error", func() {
				cmd := exec.Command(binaryPath, "--kubeconfig=/non-existing/path/kubeconfig", "get", "nginx")
				session, err := cmd.CombinedOutput()
				Expect(err).To(HaveOccurred())
				Expect(string(session)).To(ContainSubstring("failed to create config from kubeconfig path"))
			})
		})
	})

	Context("installed profiles and catalog profiles with get", func() {
		var (
			namespace        = "default"
			installationName = "long-name-to-ensure-padding"
			ctx              = context.TODO()
			pSub             profilesv1.ProfileInstallation
		)

		BeforeEach(func() {
			profileURL := "https://github.com/weaveworks/profiles-examples"
			pSub = profilesv1.ProfileInstallation{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileInstallation",
					APIVersion: "profileinstallations.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      installationName,
					Namespace: namespace,
				},
				Spec: profilesv1.ProfileInstallationSpec{
					Source: &profilesv1.Source{
						URL: profileURL,
						Tag: "weaveworks-nginx/v0.1.0",
					},
					Catalog: &profilesv1.Catalog{
						Catalog: "nginx-catalog",
						Profile: "weaveworks-nginx",
						Version: "v0.1.0",
					},
				},
			}
			Expect(kClient.Create(ctx, &pSub)).Should(Succeed())
		})

		AfterEach(func() {
			Expect(kClient.Delete(ctx, &pSub)).Should(Succeed())
		})

		It("returns all the installations and catalog profiles", func() {
			cmd := exec.Command(binaryPath, "get")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())

			expected := "PACKAGE CATALOG\n" +
				"CATALOG/PROFILE               	VERSION	DESCRIPTION                     \n" +
				"nginx-catalog/weaveworks-nginx	v0.1.0 	This installs nginx.           \t\n" +
				"nginx-catalog/weaveworks-nginx	v0.1.1 	This installs nginx.           \t\n" +
				"nginx-catalog/bitnami-nginx   	v0.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/nginx           	v2.0.1 	This installs nginx.           \t\n" +
				"nginx-catalog/some-other-nginx	       	This installs some other nginx.\t\n\n"
			Expect(string(session)).To(ContainSubstring(expected))

			expected = "INSTALLED PACKAGES\n" +
				"NAMESPACE\tNAME                       \tSOURCE                               \tAVAILABLE UPDATES \n" +
				"default  \tlong-name-to-ensure-padding\tnginx-catalog/weaveworks-nginx/v0.1.0\tv0.1.1           \t\n"
			Expect(string(session)).To(ContainSubstring(expected))
		})

		It("returns all the installations and catalog profiles with matching name", func() {
			profileURL := "https://github.com/weaveworks/profiles-examples"
			bitnamiSub := profilesv1.ProfileInstallation{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileInstallation",
					APIVersion: "profileinstallations.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bitnami-nginx",
					Namespace: namespace,
				},
				Spec: profilesv1.ProfileInstallationSpec{
					Source: &profilesv1.Source{
						URL: profileURL,
						Tag: "bitnami-nginx/v0.0.1",
					},
					Catalog: &profilesv1.Catalog{
						Catalog: "nginx-catalog",
						Profile: "bitnami-nginx",
						Version: "v0.0.1",
					},
				},
			}
			Expect(kClient.Create(ctx, &bitnamiSub)).Should(Succeed())

			cmd := exec.Command(binaryPath, "get", "bitnami-nginx")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())

			expected := "INSTALLED PACKAGES\n" +
				"NAMESPACE\tNAME         \tSOURCE                            \tAVAILABLE UPDATES \n" +
				"default  \tbitnami-nginx\tnginx-catalog/bitnami-nginx/v0.0.1\t-                \t\n\n" +
				"PACKAGE CATALOG\n" +
				"CATALOG/PROFILE            	VERSION	DESCRIPTION          \n" +
				"nginx-catalog/bitnami-nginx	v0.0.1 	This installs nginx.\t\n\n"
			Expect(string(session)).To(ContainSubstring(expected))
		})
	})

	Context("installed profiles with get", func() {
		var (
			namespace        = "default"
			installationName = "long-name-to-ensure-padding"
			ctx              = context.TODO()
			pSub             profilesv1.ProfileInstallation
		)

		BeforeEach(func() {
			profileURL := "https://github.com/weaveworks/profiles-examples"
			pSub = profilesv1.ProfileInstallation{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileInstallation",
					APIVersion: "profileinstallations.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      installationName,
					Namespace: namespace,
				},
				Spec: profilesv1.ProfileInstallationSpec{
					Source: &profilesv1.Source{
						URL: profileURL,
						Tag: "weaveworks-nginx/v0.1.0",
					},
					Catalog: &profilesv1.Catalog{
						Catalog: "nginx-catalog",
						Profile: "weaveworks-nginx",
						Version: "v0.1.0",
					},
				},
			}
			Expect(kClient.Create(ctx, &pSub)).Should(Succeed())
		})

		AfterEach(func() {
			Expect(kClient.Delete(ctx, &pSub)).Should(Succeed())
		})

		It("returns the installations", func() {
			getCmd := func() []string {
				cmd := exec.Command(binaryPath, "get", "--installed")
				session, err := cmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				return strings.Split(string(session), "\n")
			}

			Eventually(getCmd).Should(ContainElements(
				"INSTALLED PACKAGES",
				"NAMESPACE\tNAME                       \tSOURCE                               \tAVAILABLE UPDATES ",
				"default  \tlong-name-to-ensure-padding\tnginx-catalog/weaveworks-nginx/v0.1.0\tv0.1.1           \t",
			))
		})

		It("doesnt not break json output", func() {
			cmd := exec.Command(binaryPath, "get", "-o", "json", "--installed")
			session, err := cmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(session)).ToNot(ContainSubstring("INSTALLED PACKAGES"))
		})

		When("there are no available updates", func() {
			It("returns the installations", func() {
				profileURL := "https://github.com/weaveworks/profiles-examples"
				bitnamiSub := profilesv1.ProfileInstallation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileInstallation",
						APIVersion: "profileinstallations.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bitnami-profile",
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileInstallationSpec{
						Source: &profilesv1.Source{
							URL: profileURL,
							Tag: "bitnami-nginx/v0.0.1",
						},
						Catalog: &profilesv1.Catalog{
							Catalog: "nginx-catalog",
							Profile: "bitnami-nginx",
							Version: "v0.0.1",
						},
					},
				}
				Expect(kClient.Create(ctx, &bitnamiSub)).Should(Succeed())
				getCmd := func() []string {
					cmd := exec.Command(binaryPath, "get", "--installed")
					session, err := cmd.CombinedOutput()
					Expect(err).ToNot(HaveOccurred())
					return strings.Split(string(session), "\n")
				}

				Eventually(getCmd).Should(ContainElements(
					"INSTALLED PACKAGES",
					"NAMESPACE\tNAME                       \tSOURCE                               \tAVAILABLE UPDATES ",
					"default  \tbitnami-profile            \tnginx-catalog/bitnami-nginx/v0.0.1   \t-                \t",
					"default  \tlong-name-to-ensure-padding\tnginx-catalog/weaveworks-nginx/v0.1.0\tv0.1.1           \t",
				))
				Expect(kClient.Delete(ctx, &bitnamiSub)).Should(Succeed())
			})

			It("returns the installation matching name", func() {
				profileURL := "https://github.com/weaveworks/profiles-examples"
				bitnamiSub := profilesv1.ProfileInstallation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ProfileInstallation",
						APIVersion: "profileinstallations.weave.works/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bitnami-profile",
						Namespace: namespace,
					},
					Spec: profilesv1.ProfileInstallationSpec{
						Source: &profilesv1.Source{
							URL: profileURL,
							Tag: "bitnami-nginx/v0.0.1",
						},
						Catalog: &profilesv1.Catalog{
							Catalog: "nginx-catalog",
							Profile: "bitnami-nginx",
							Version: "v0.0.1",
						},
					},
				}
				Expect(kClient.Create(ctx, &bitnamiSub)).Should(Succeed())
				getCmd := func() []string {
					cmd := exec.Command(binaryPath, "get", "--installed", "bitnami-profile")
					session, err := cmd.CombinedOutput()
					Expect(err).ToNot(HaveOccurred())
					return strings.Split(string(session), "\n")
				}

				Eventually(getCmd).Should(ConsistOf(
					"INSTALLED PACKAGES",
					"NAMESPACE\tNAME           \tSOURCE                            \tAVAILABLE UPDATES ",
					"default  \tbitnami-profile\tnginx-catalog/bitnami-nginx/v0.0.1\t-                \t",
					"",
					"",
				))
				Expect(kClient.Delete(ctx, &bitnamiSub)).Should(Succeed())
			})
		})
	})

	Context("version", func() {
		When("version is used in the get command", func() {
			It("shows the right profile", func() {
				cmd := exec.Command(binaryPath, "get", "--profile-version", "v0.1.0", "nginx-catalog/weaveworks-nginx")
				output, err := cmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(string(output)).To(ContainSubstring("Catalog      \tnginx-catalog                                      \t\n" +
					"Name         \tweaveworks-nginx                                   \t\n" +
					"Version      \tv0.1.0                                             \t\n" +
					"Description  \tThis installs nginx.                               \t\n" +
					"URL          \thttps://github.com/weaveworks/profiles-examples    \t\n" +
					"Maintainer   \tweaveworks (https://github.com/weaveworks/profiles)\t\n" +
					"Prerequisites\tKubernetes 1.18+                                   \t\n",
				))
			})
		})

		When("-o is set to json", func() {
			It("returns the profile info in json", func() {
				cmd := exec.Command(binaryPath, "get", "-o", "json", "--profile-version", "v0.1.0", "nginx-catalog/weaveworks-nginx")
				session, err := cmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(string(session)).To(ContainSubstring(`{
  "tag": "weaveworks-nginx/v0.1.0",
  "catalogSource": "nginx-catalog",
  "url": "https://github.com/weaveworks/profiles-examples",
  "name": "weaveworks-nginx",
  "description": "This installs nginx.",
  "maintainer": "weaveworks (https://github.com/weaveworks/profiles)",
  "prerequisites": [
    "Kubernetes 1.18+"
  ]
}`))
			})
		})

		When("a name argument is not provided correctly", func() {
			It("returns a useful error", func() {
				cmd := exec.Command(binaryPath, "get", "--profile-version", "v0.1.0", "test-profile")
				session, err := cmd.CombinedOutput()
				Expect(err).To(HaveOccurred())
				Expect(string(session)).To(ContainSubstring("both catalog name and profile name must be provided example: pctl get catalog/weaveworks-nginx --version v0.1.0"))
			})
		})
	})
})
