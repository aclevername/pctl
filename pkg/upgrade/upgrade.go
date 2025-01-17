package upgrade

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	copypkg "github.com/otiai10/copy"
	"github.com/weaveworks/pctl/pkg/catalog"
	"github.com/weaveworks/pctl/pkg/git"
	"github.com/weaveworks/pctl/pkg/install"
	"github.com/weaveworks/pctl/pkg/runner"
	"github.com/weaveworks/pctl/pkg/upgrade/repo"
	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"
	"sigs.k8s.io/yaml"
)

// UpgradeCfg holds the fields used during upgrades a installation
type UpgradeConfig struct {
	ProfileDir     string
	Version        string
	CatalogClient  catalog.CatalogClient
	CatalogManager catalog.CatalogManager
	RepoManager    repo.RepoManager
	WorkingDir     string
	Message        string
}

var copy func(src, dest string) error = func(src, dest string) error {
	return copypkg.Copy(src, dest)
}

// Upgrade the profiel installation to a new version
func Upgrade(cfg UpgradeConfig) error {
	out, err := ioutil.ReadFile(path.Join(cfg.ProfileDir, "profile-installation.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read profile installation: %w", err)
	}

	var profileInstallation profilesv1.ProfileInstallation
	if err := yaml.Unmarshal(out, &profileInstallation); err != nil {
		return fmt.Errorf("failed to parse profile installation: %w", err)
	}

	fmt.Printf("upgrading profile %q from version %q to %q\n", profileInstallation.Name, profileInstallation.Spec.Catalog.Version, cfg.Version)

	var gitRepoName, gitRepoNamespace string
	catalogName := profileInstallation.Spec.Catalog.Catalog
	profileName := profileInstallation.Spec.Catalog.Profile
	currentVersion := profileInstallation.Spec.Catalog.Version
	if profileInstallation.Spec.GitRepository != nil {
		gitRepoName = profileInstallation.Spec.GitRepository.Name
		gitRepoNamespace = profileInstallation.Spec.GitRepository.Namespace
	}

	//check new version exists
	_, err = cfg.CatalogManager.Show(cfg.CatalogClient, catalogName, profileName, cfg.Version)
	if err != nil {
		return fmt.Errorf("failed to get profile %q in catalog %q version %q: %w", profileName, catalogName, cfg.Version, err)
	}

	err = cfg.RepoManager.CreateRepoWithContent(func() error {
		installConfig := catalog.InstallConfig{
			Clients: catalog.Clients{
				CatalogClient: cfg.CatalogClient,
				Installer: install.NewInstaller(install.Config{
					GitClient: git.NewCLIGit(git.CLIGitConfig{
						Message: cfg.Message,
					}, &runner.CLIRunner{}),
					RootDir:          cfg.WorkingDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				}),
			},
			Profile: catalog.Profile{
				ProfileConfig: catalog.ProfileConfig{
					ProfileName: profileName,
					CatalogName: catalogName,
					Version:     currentVersion,
					ConfigMap:   profileInstallation.Spec.ConfigMap,
				},
				GitRepoConfig: catalog.GitRepoConfig{
					Namespace: gitRepoNamespace,
					Name:      gitRepoName,
				},
			},
		}
		if err := cfg.CatalogManager.Install(installConfig); err != nil {
			return fmt.Errorf("failed to install base profile: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create repository for upgrade: %w", err)
	}

	err = cfg.RepoManager.CreateBranchWithContentFromMain("user-changes", func() error {
		if err := copy(cfg.ProfileDir, cfg.WorkingDir); err != nil {
			return fmt.Errorf("failed to copy profile during upgrade: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create branch with user changes: %w", err)
	}

	err = cfg.RepoManager.CreateBranchWithContentFromMain("update-changes", func() error {
		installConfig := catalog.InstallConfig{
			Clients: catalog.Clients{
				CatalogClient: cfg.CatalogClient,
				Installer: install.NewInstaller(install.Config{
					GitClient: git.NewCLIGit(git.CLIGitConfig{
						Message: cfg.Message,
					}, &runner.CLIRunner{}),
					RootDir:          cfg.WorkingDir,
					GitRepoNamespace: gitRepoNamespace,
					GitRepoName:      gitRepoName,
				}),
			},
			Profile: catalog.Profile{
				ProfileConfig: catalog.ProfileConfig{
					ProfileName: profileName,
					CatalogName: catalogName,
					Version:     cfg.Version,
					ConfigMap:   profileInstallation.Spec.ConfigMap,
				},
				GitRepoConfig: catalog.GitRepoConfig{
					Namespace: gitRepoNamespace,
					Name:      gitRepoName,
				},
			},
		}

		if err := cfg.CatalogManager.Install(installConfig); err != nil {
			return fmt.Errorf("failed to install update profile: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create branch with update changes: %w", err)
	}

	mergeConflictFiles, err := cfg.RepoManager.MergeBranches("update-changes", "user-changes")
	if err != nil {
		return fmt.Errorf("failed to merge updates with user changes: %w", err)
	}

	if err := os.RemoveAll(cfg.ProfileDir); err != nil {
		return fmt.Errorf("failed to remove existing profile installation: %w", err)
	}

	if err := os.RemoveAll(filepath.Join(cfg.WorkingDir, ".git/")); err != nil {
		return fmt.Errorf("failed to remove git directory from upgrade directory: %w", err)
	}

	if err := copy(cfg.WorkingDir, cfg.ProfileDir); err != nil {
		return fmt.Errorf("failed to copy upgraded installation into installation directory: %w", err)
	}

	if len(mergeConflictFiles) > 0 {
		msg := "upgrade succeeded but merge conflicts have occurred, please resolve manually. Files containing conflicts:\n"
		for _, mergeConflictFile := range mergeConflictFiles {
			msg = fmt.Sprintf("%s- %s\n", msg, filepath.Join(cfg.ProfileDir, mergeConflictFile))
		}
		msg = strings.TrimSuffix(msg, "\n")
		return fmt.Errorf(msg)
	}

	fmt.Println("upgrade completed successfully")
	return nil
}
