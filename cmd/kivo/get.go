package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	profilesv1 "github.com/weaveworks/profiles/api/v1alpha1"

	"github.com/weaveworks/kivo-cli/pkg/catalog"
	"github.com/weaveworks/kivo-cli/pkg/client"
	"github.com/weaveworks/kivo-cli/pkg/formatter"
	"github.com/weaveworks/kivo-cli/pkg/log"
)

func getCmd() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get a profile",
		UsageText: "kivo get [--output table|json <QUERY> --installed --catalog --version] \n\n" +
			"   example: kivo get nginx",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				DefaultText: "table",
				Value:       "table",
				Usage:       "Output format. json|table",
			},
			&cli.BoolFlag{
				Name:    "installed",
				Aliases: []string{"i"},
				Usage:   "Get all installed profiles",
			},
			&cli.BoolFlag{
				Name:    "catalog",
				Aliases: []string{"c"},
				Usage:   "Get all profiles in a catalog",
			},
			&cli.StringFlag{
				Name:    "profile-version",
				Aliases: []string{"p"},
				Usage: "Get information about a profile \n\n" +
					"   kivo [--kubeconfig=<kubeconfig-path>] get <CATALOG>/<PROFILE>\n\n" +
					"   example: kivo get catalog/weaveworks-nginx --profile-version v0.1.0",
			},
		},
		Action: func(c *cli.Context) error {
			outFormat := c.String("output")

			cl, err := buildK8sClient(c.String("kubeconfig"))
			if err != nil {
				return err
			}
			catalogClient, err := getCatalogClient(c)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "get")
				return err
			}

			// get all installed and catalog profiles
			if c.Args().Len() < 1 {
				if c.Bool("installed") {
					return getInstalledProfiles(cl, catalogClient, "", outFormat)
				}
				if c.Bool("catalog") {
					return getCatalogProfiles(catalogClient, "", outFormat)
				}

				err := getInstalledProfiles(cl, catalogClient, "", outFormat)
				if err != nil {
					return err
				}
				return getCatalogProfiles(catalogClient, "", outFormat)
			} else {
				name, catalogClient, err := parseArgs(c)
				if err != nil {
					_ = cli.ShowCommandHelp(c, "get")
					return err
				}

				// check if other flags are passed along with the name
				if c.String("profile-version") != "" {
					return getCatalogProfilesWithVersion(c, catalogClient, name, c.String("profile-version"), outFormat)
				}

				if c.Bool("catalog") {
					return getCatalogProfiles(catalogClient, name, outFormat)
				}

				if c.Bool("installed") {
					return getInstalledProfiles(cl, catalogClient, name, outFormat)
				}

				// default return all installed and catalog profiles with name
				err = getInstalledProfiles(cl, catalogClient, name, outFormat)
				if err != nil {
					return err
				}

				err = getCatalogProfiles(catalogClient, name, outFormat)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func getCatalogProfiles(catalogClient *client.Client, name string, outFormat string) error {
	manager := &catalog.Manager{}
	profiles, err := manager.Search(catalogClient, name)
	if err != nil {
		return err
	}
	return formatOutput(profiles, outFormat)
}

func getInstalledProfiles(cl runtimeclient.Client, catalogClient *client.Client, name string, outFormat string) error {
	manager := &catalog.Manager{}
	data, err := manager.List(cl, catalogClient, name)
	if err != nil {
		return err
	}

	return formatInstalledProfilesOutput(data, outFormat)
}

func getCatalogProfilesWithVersion(c *cli.Context, catalogClient *client.Client, name string, version string, outFormat string) error {
	parts := strings.Split(name, "/")
	if len(parts) < 2 {
		_ = cli.ShowCommandHelp(c, "get")
		return fmt.Errorf("both catalog name and profile name must be provided example: kivo get catalog/weaveworks-nginx --version v0.1.0")
	}
	catalogName, profileName := parts[0], parts[1]
	manager := &catalog.Manager{}
	profile, err := manager.Show(catalogClient, catalogName, profileName, version)
	if err != nil {
		return err
	}

	return formatCatlogProfilesOutput(profile, outFormat)
}

func profilesDataFunc(profiles []profilesv1.ProfileCatalogEntry) func() interface{} {
	return func() interface{} {
		tc := formatter.TableContents{
			Headers: []string{"Catalog/Profile", "Version", "Description"},
		}
		for _, profile := range profiles {
			tc.Data = append(tc.Data, []string{
				fmt.Sprintf("%s/%s", profile.CatalogSource, profile.Name),
				profilesv1.GetVersionFromTag(profile.Tag),
				profile.Description,
			})
		}
		return tc
	}
}

func installedProfilesDataFunc(data []catalog.ProfileData) func() interface{} {
	return func() interface{} {
		tc := formatter.TableContents{
			Headers: []string{"Namespace", "Name", "Source", "Available Updates"},
		}
		for _, d := range data {
			source := fmt.Sprintf("%s/%s/%s", d.Profile.Catalog, d.Profile.Profile, d.Profile.Version)
			if d.Profile.Catalog == "-" {
				source = fmt.Sprintf("%s:%s:%s", d.Profile.URL, d.Profile.Branch, d.Profile.Path)
			}
			tc.Data = append(tc.Data, []string{
				d.Profile.Namespace,
				d.Profile.Name,
				source,
				strings.Join(d.AvailableVersionUpdates, ","),
			})
		}
		return tc
	}
}

func profileWithVersionDataFunc(profile profilesv1.ProfileCatalogEntry) func() interface{} {
	return func() interface{} {
		return formatter.TableContents{
			Data: [][]string{
				{"Catalog", profile.CatalogSource},
				{"Name", profile.Name},
				{"Version", profilesv1.GetVersionFromTag(profile.Tag)},
				{"Description", profile.Description},
				{"URL", profile.URL},
				{"Maintainer", profile.Maintainer},
				{"Prerequisites", strings.Join(profile.Prerequisites, ", ")},
			},
		}
	}
}

func formatOutput(profiles []profilesv1.ProfileCatalogEntry, outFormat string) error {
	if len(profiles) == 0 {
		log.Failuref("No available catalog profiles found.")
		return nil
	}

	var f formatter.Formatter
	f = formatter.NewTableFormatter()
	getter := profilesDataFunc(profiles)

	if outFormat == "json" {
		f = formatter.NewJSONFormatter()
		getter = func() interface{} { return profiles }
	}

	out, err := f.Format(getter)
	if err != nil {
		return err
	}

	if outFormat == "table" {
		fmt.Println("PACKAGE CATALOG")
	}

	fmt.Println(out)
	return nil
}

func formatInstalledProfilesOutput(data []catalog.ProfileData, outFormat string) error {
	if len(data) == 0 {
		log.Failuref("No installed profiles found.")
		return nil
	}

	var f formatter.Formatter
	f = formatter.NewTableFormatter()
	getter := installedProfilesDataFunc(data)

	if outFormat == "json" {
		f = formatter.NewJSONFormatter()
		getter = func() interface{} { return data }
	}

	out, err := f.Format(getter)
	if err != nil {
		return err
	}

	if outFormat == "table" {
		fmt.Println("INSTALLED PACKAGES")
	}

	fmt.Println(out)
	return nil
}

func formatCatlogProfilesOutput(profile profilesv1.ProfileCatalogEntry, outFormat string) error {
	if profile.Name == "" {
		log.Failuref("No profile found.")
		return nil
	}

	var f formatter.Formatter
	f = formatter.NewTableFormatter()
	getter := profileWithVersionDataFunc(profile)

	if outFormat == "json" {
		f = formatter.NewJSONFormatter()
		getter = func() interface{} { return profile }
	}

	out, err := f.Format(getter)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
