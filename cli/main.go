package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thaddock/modrinth-query/api"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{

		Name:  "Modrinth Query",
		Usage: "TBD",
		Commands: []*cli.Command{
			{
				Name:  "project",
				Usage: "Project subcommands",
				Subcommands: []*cli.Command{
					{
						Name:  "query",
						Usage: "Query projects",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "query",
								Value: "",
								Usage: "query to perform",
							},
							&cli.StringFlag{
								Name:  "facets",
								Value: "",
								Usage: "query facets",
							},
						},
						Action: func(ctx *cli.Context) error {
							fmt.Println("*** Query projects")
							c := api.NewClient()
							sr, err := api.ProjectSearch(c, ctx.String("query"), ctx.String("facets"), "", 0, 10)
							if err != nil {
								return err
							}
							fmt.Printf("TotalHits: %d\n", sr.TotalHits)
							for _, project := range sr.Hits {
								fmt.Printf("%s: %s\n", project.Slug, project.Description)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "version",
				Usage: "Version subcommands",
				Subcommands: []*cli.Command{
					{
						Name:  "query",
						Usage: "Get version",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id",
								Value: "",
								Usage: "Project id",
							},
							&cli.StringFlag{
								Name:  "loaders",
								Value: "",
								Usage: "Loaders",
							},
							&cli.StringSliceFlag{
								Name:  "game_versions",
								Value: nil,
								Usage: "Game versions",
							},
							&cli.StringSliceFlag{
								Name:  "featured",
								Value: nil,
								Usage: "Featured?",
							},
						},
						Action: func(ctx *cli.Context) error {
							c := api.NewClient()
							var featured *bool
							if ctx.String("featured") != "" {
								value := ctx.Bool("featured")
								featured = &value
							}
							vr, err := api.GetVersion(c, ctx.String("id"), ctx.StringSlice("loaders"), ctx.StringSlice("game_versions"), featured)
							if err != nil {
								return err
							}
							for _, v := range vr {
								fmt.Printf("%s:\t%s\t%s\t%s\n", v.Name, v.VersionNumber, v.DatePublished, v.GameVersions)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "version_file",
				Usage: "Version File subcommands",
				Subcommands: []*cli.Command{
					{
						Name:  "query",
						Usage: "Get Version from File",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "hash",
								Value: "",
								Usage: "Hash to query",
							},
							&cli.StringFlag{
								Name:  "algorithm",
								Value: "sha1",
								Usage: "Algorithm for hash",
							},
						},
						Action: func(ctx *cli.Context) error {
							c := api.NewClient()
							vr, err := api.GetVersionFile(c, ctx.String("hash"), ctx.String("algorithm"), false)
							if err != nil {
								return err
							}
							if vr == nil {
								fmt.Printf("File not found\n")
								return nil
							}
							fmt.Printf("%s %s\n", vr.Name, vr.VersionNumber)
							return nil
						},
					},
				},
			},
			{
				Name:  "sync",
				Usage: "Sync a directory",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "dir",
						Value: "",
						Usage: "Directory to scan",
					},
					&cli.StringSliceFlag{
						Name:  "loaders",
						Value: nil,
						Usage: "Loaders",
					},
					&cli.StringSliceFlag{
						Name:  "game_versions",
						Value: nil,
						Usage: "Game versions",
					},
					&cli.BoolFlag{
						Name:  "update",
						Value: false,
						Usage: "Update modules with newer versions",
					},
				},
				Action: func(ctx *cli.Context) error {
					c := api.NewClient()
					return DoSync(c, ctx.String("dir"), ctx.StringSlice("loaders"), ctx.StringSlice("game_versions"), ctx.Bool("update"))
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
