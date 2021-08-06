package cluster

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"

	"github.com/weaveworks/pctl/pkg/runner"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	kubectlCmd = "kubectl"
	// profiles bundles ready to be installed files under `install`. The rest of the resources
	// are left for manual configuration.
	installManifestFile = "prepare.yaml"
	namespace           = "profiles-system"
)

// FluxCRDs are CRDs which install is checking if they are present in the cluster or not.
var FluxCRDs = []string{
	strings.ToLower(helmv2.HelmReleaseKind),
	strings.ToLower(kustomizev1.KustomizationKind),
	strings.ToLower(sourcev1.BucketKind),
	strings.ToLower(sourcev1.GitRepositoryKind),
	strings.ToLower(sourcev1.HelmChartKind),
	strings.ToLower(sourcev1.HelmRepositoryKind),
}

// Fetcher will download a manifest tar file from a remote repository.
type Fetcher struct {
	Client *http.Client
}

// Applier applies the previously generated manifest files.
type Applier struct {
	Runner runner.Runner
	Waiter Waiter
}

// Installer will install an environment.
type Installer struct {
	InstallConfig
	Applier *Applier
	Fetcher *Fetcher
	Runner  runner.Runner
}

// InstallConfig defines configuration options for install.
type InstallConfig struct {
	// BaseURL is given even one would like to download manifests from a fork
	// or a test repo.
	BaseURL               string
	Location              string
	Version               string
	KubeContext           string
	KubeConfig            string
	FluxNamespace         string
	IgnorePreflightErrors bool
	DryRun                bool
	Keep                  bool
	K8sClient             client.Client
}

// NewInstaller creates an installer with set dependencies ready to be used.
func NewInstaller(cfg InstallConfig) (*Installer, error) {
	if cfg.Location == "" {
		tmp, err := ioutil.TempDir("", "pctl-manifests")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp folder for manifest files: %w", err)
		}
		cfg.Location = tmp
	}
	r := &runner.CLIRunner{}
	return &Installer{
		InstallConfig: cfg,
		Fetcher: &Fetcher{
			Client: http.DefaultClient,
		},
		Applier: &Applier{
			Waiter: NewKubeWaiter(KubeConfig{
				Client:    cfg.K8sClient,
				Interval:  5 * time.Second,
				Timeout:   15 * time.Minute,
				Namespace: namespace,
			}),
			Runner: r,
		},
		Runner: r,
	}, nil
}

// Install will install an environment with everything that is needed to run profiles.
func (i *Installer) Install() error {
	defer func() {
		if i.Keep {
			return
		}
		if err := os.RemoveAll(i.Location); err != nil {
			fmt.Printf("failed to remove temporary folder at location: %s. Please clean manually.", i.Location)
		}
	}()
	if err := i.PreFlightCheck(); err != nil {
		return err
	}
	if err := i.Fetcher.Fetch(context.Background(), i.BaseURL, i.Version, i.Location); err != nil {
		return err
	}
	return i.Applier.Apply(i.Location, i.KubeContext, i.KubeConfig, i.DryRun)
}

// PreFlightCheck checks whether install can run or not.
func (i *Installer) PreFlightCheck() error {
	fmt.Print("Checking if flux namespace exists...")
	args := []string{"get", "namespace", i.FluxNamespace, "--output", "name"}
	if output, err := i.Runner.Run(kubectlCmd, args...); err != nil {
		fmt.Println("\nOutput from kubectl command: ", string(output))
		if i.IgnorePreflightErrors {
			fmt.Println("WARNING: failed to get flux namespace. Flux is required for profiles to work.")
		} else {
			return fmt.Errorf("failed to get flux namespace: %w\nTo ignore this error, please see the  --ignore-preflight-checks flag.", err)
		}
	}
	fmt.Println("done.")
	fmt.Print("Checking for flux CRDs...")
	output, err := i.Runner.Run(kubectlCmd, "get", "crds", "--output", "jsonpath='{.items[*].spec.names.singular}'")
	if err != nil {
		if i.IgnorePreflightErrors {
			fmt.Println("WARNING: failed to list all installed crds. Flux is required for profiles to work.")
		} else {
			return fmt.Errorf("failed to list all installed crds: %w", err)
		}
	}
	// the output contains an opening an closing '
	output = bytes.Trim(output, "'")
	// create an easily searchable list of installed CRDs for verification
	crds := map[string]struct{}{}
	for _, c := range strings.Split(string(output), " ") {
		crds[c] = struct{}{}
	}
	for _, crd := range FluxCRDs {
		if _, ok := crds[crd]; !ok {
			if i.IgnorePreflightErrors {
				fmt.Println("WARNING: failed to find flux crd resource. Flux is required for profiles to work.")
			} else {
				return fmt.Errorf("failed to get crd %s\nTo ignore this error, please see the  --ignore-preflight-checks flag.", crd)
			}
		}
	}

	fmt.Println("done.")
	return nil
}

// Fetch the latest or a version of the released manifest files for profiles.
func (f *Fetcher) Fetch(ctx context.Context, url, version, dir string) error {
	ghURL := fmt.Sprintf("%s/latest/download/%s", url, installManifestFile)
	hasVersionPrefix := strings.HasPrefix(version, "v")
	if hasVersionPrefix {
		ghURL = fmt.Sprintf("%s/download/%s/%s", url, version, installManifestFile)
	}

	req, err := http.NewRequest("GET", ghURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for %s, error: %w", ghURL, err)
	}

	resp, err := f.Client.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to download prepare.yaml from %s, error: %w", ghURL, err)
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Println("Failed to close body reader.")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download prepare.yaml from %s, status: %s", ghURL, resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body of the response: %w", err)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, installManifestFile), content, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write out file to location: %w", err)
	}

	return nil
}

// Apply applies the fetched manifest files to a cluster.
func (a *Applier) Apply(folder string, kubeContext string, kubeConfig string, dryRun bool) error {
	kubectlArgs := []string{"apply", "-f", filepath.Join(folder, installManifestFile)}
	if dryRun {
		kubectlArgs = append(kubectlArgs, "--dry-run=client", "--output=yaml")
	}
	if kubeContext != "" {
		kubectlArgs = append(kubectlArgs, "--context="+kubeContext)
	}
	if kubeConfig != "" {
		kubectlArgs = append(kubectlArgs, "--kubeconfig="+kubeConfig)
	}
	output, err := a.Runner.Run(kubectlCmd, kubectlArgs...)
	if err != nil {
		fmt.Println("Log from kubectl: ", string(output))
		return fmt.Errorf("install failed: %w", err)
	}
	if dryRun {
		fmt.Print(string(output))
		return nil
	}
	fmt.Print("Waiting for resources to be ready...")
	if err := a.Waiter.Wait("profiles-controller-manager"); err != nil {
		return fmt.Errorf("failed to wait for resources to be ready: %w", err)
	}
	fmt.Println("done.")
	return nil
}
