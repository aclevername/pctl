package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/weaveworks/pctl/pkg/cluster"
	"github.com/weaveworks/pctl/pkg/config"
)

const (
	fluxNamespace = "flux-system"
)

func prepareCmd() *cli.Command {
	return &cli.Command{
		Name:      "prepare",
		Usage:     "prepares the cluster for profiles by deploying the profile controllers and custom resource definitions",
		UsageText: "pctl prepare",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "If defined, nothing will be applied.",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "keep",
				Usage: "Keep the downloaded manifest files.",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "ignore-preflight-errors",
				Usage: "Instead of stopping the process, output warnings when they occur during preflight check.",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "version",
				Usage: "Define the tagged version to use which can be found under releases in the profiles repository. Exp: [v]0.0.1",
			},
			&cli.StringFlag{
				Name:        "baseurl",
				Usage:       "Define the url to go and fetch releases from.",
				Value:       releaseUrl,
				DefaultText: releaseUrl,
			},
			&cli.StringFlag{
				Name:        "flux-namespace",
				Usage:       "Define the namespace in which flux is installed.",
				Value:       fluxNamespace,
				DefaultText: fluxNamespace,
			},
			&cli.StringFlag{
				Name:        "out",
				Usage:       "Specify the output location of the downloaded prepare file.",
				Value:       "",
				DefaultText: "os.Temp",
			},
			&cli.StringFlag{
				Name:  "context",
				Usage: "The Kubernetes context to use to apply the manifest files .",
			},
			&cli.StringFlag{
				Name:  "git-repository",
				Value: "",
				Usage: "The namespace and name of the GitRepository object governing the flux repo.",
			},
		},
		Action: func(c *cli.Context) error {
			kubeConfig := c.String("kubeconfig")
			k8sClient, err := buildK8sClient(kubeConfig)
			if err != nil {
				return fmt.Errorf("failed to build kubernetes client: %w", err)
			}
			p, err := cluster.NewPreparer(cluster.PrepConfig{
				BaseURL:               c.String("baseurl"),
				Version:               c.String("version"),
				KubeConfig:            c.String("kubeconfig"),
				KubeContext:           c.String("context"),
				FluxNamespace:         c.String("flux-namespace"),
				Location:              c.String("out"),
				DryRun:                c.Bool("dry-run"),
				Keep:                  c.Bool("keep"),
				IgnorePreflightErrors: c.Bool("ignore-preflight-errors"),
				K8sClient:             k8sClient,
			})
			if err != nil {
				return err
			}

			if err := p.Prepare(); err != nil {
				return err
			}

			return config.Create(c.String("git-repository"))
		},
	}
}
