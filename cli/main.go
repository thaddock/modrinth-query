package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lazerdye/modrinth_query/api"
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
						},
						Action: func(ctx *cli.Context) error {
							fmt.Println("*** Query projects")
							c := api.NewClient()
							sr, err := api.ProjectSearch(c, ctx.String("query"), "", "", 0, 10)
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
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
