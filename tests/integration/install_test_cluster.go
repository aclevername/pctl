package integration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	waitForDeploymentArgs = []string{
		"-n",
		"profiles-system",
		"wait",
		"--for=condition=available",
		"deployment",
		"profiles-controller-manager",
		"--timeout",
		"5m",
	}
	waitForPodsArgs = []string{
		"-n",
		"profiles-system",
		"wait",
		"--for=condition=Ready",
		"--all",
		"pods",
		"--timeout",
		"5m",
	}
	applySourceCatalogArgs = []string{
		"apply",
		"-f",
		"catalog-source.yaml",
	}
)

// installTestCluster will create a test cluster using kivo `install` command.
// @binary -- location of the built kivo binary.
func InstallClusterComponents(binaryPath string) error {
	tmp, err := ioutil.TempDir("", "install_integration_test_01")
	if err != nil {
		return fmt.Errorf("failed to create temp folder for test: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmp); err != nil {
			fmt.Printf("failed to remove temporary folder at location: %s. Please clean manually.", tmp)
		}
	}()
	cmd := exec.Command(binaryPath, "install", "--dry-run", "--out", tmp, "--keep")
	output, err := cmd.CombinedOutput()
	if err != nil || !bytes.Contains(output, []byte("kind: List")) {
		fmt.Println("Output of install was: ", string(output))
		return fmt.Errorf("failed to run install command: %w", err)
	}
	fmt.Println("Install file generated successfully.")

	content, err := ioutil.ReadFile(filepath.Join(tmp, "prepare.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read prepare.yaml from location %s: %w", tmp, err)
	}
	fmt.Print("Replacing controller image to localhost:5000...")
	re := regexp.MustCompile(`weaveworks/profiles-controller:.*`)
	out := re.ReplaceAllString(string(content), "localhost:5000/profiles-controller:latest")
	fmt.Println("done.")

	fmt.Print("Applying modified prepare.yaml...")
	applyInstallArgs := []string{"apply", "-f", "-"}
	cmd = exec.Command("kubectl", applyInstallArgs...)
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get in pipe for kubectl: %w", err)
	}

	go func() {
		defer func(in io.WriteCloser) {
			_ = in.Close()
		}(in)
		if _, err := io.WriteString(in, out); err != nil {
			fmt.Println("Failed to write to kubectl apply: ", err)
		}
	}()

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("\nOutput from kubectl apply: ", string(output))
		return fmt.Errorf("failed to apply prepare.yaml: %w", err)
	}
	fmt.Println("done.")

	fmt.Print("Waiting for deployment...")
	if err := runKubectl(waitForDeploymentArgs); err != nil {
		return fmt.Errorf("failed to wait for deployment: %w", err)
	}

	fmt.Print("Waiting for pods to be active...")
	if err := runKubectl(waitForPodsArgs); err != nil {
		return fmt.Errorf("failed to wait for pods to be active: %w", err)
	}

	fmt.Print("Applying test catalog...")
	if err := runKubectl(applySourceCatalogArgs); err != nil {
		return fmt.Errorf("failed to apply test catalog: %w", err)
	}

	fmt.Println("Happy testing!")
	return nil
}

func runKubectl(args []string) error {
	cmd := exec.Command("kubectl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error from kubectl: ", string(output))
		return err
	}
	fmt.Println("done.")
	return nil
}
