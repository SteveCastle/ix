package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SteveCastle/ix"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "ix",
		Usage:   "IX is a file indexer that uses a simple file structure and symlinks to store metadata about files.",
		Version: "0.0.1",
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "store",
				Aliases: []string{"s"},
				Usage:   "Store returns the store location that will be used for saving metadata.",
				Action: func(c *cli.Context) error {
					fmt.Println(ix.FindStore("./"))
					return nil
				},
				Subcommands: []*cli.Command{},
			},
			{
				Name:      "tag",
				Aliases:   []string{"t"},
				Usage:     "Tag a file with a category and name. If a directory is specified, all files in the directory will be tagged.",
				ArgsUsage: "<category> <name> <file>",
				BashComplete: func(cCtx *cli.Context) {
					if cCtx.NArg() < 2 {
						categories := ix.ListCategories()
						argSubstring := cCtx.Args().Get(0)

						foundCategories := []string{}
						for _, category := range categories {
							if strings.Contains(category, argSubstring) {
								foundCategories = append(foundCategories, category)
							}
						}
						if len(foundCategories) == 0 {
							for _, t := range categories {
								fmt.Println(t)
							}
							return
						}
						for _, t := range foundCategories {
							fmt.Println(t)
						}
					}
					if cCtx.NArg() == 2 || (cCtx.NArg() == 1 && strings.HasSuffix(cCtx.Args().Get(0), " ")) {
						tags := ix.ListTags(cCtx.Args().Get(0))
						argSubstring := cCtx.Args().Get(1)
						foundTags := []string{}
						for _, tag := range tags {
							if strings.Contains(tag, argSubstring) {
								foundTags = append(foundTags, tag)
							}
						}
						if len(foundTags) == 0 {
							for _, t := range tags {
								fmt.Println(t)
							}
							return
						}
						for _, t := range foundTags {
							fmt.Println(t)
						}
					}
				},
				Subcommands: []*cli.Command{
					{
						Name:      "create",
						Aliases:   []string{"c"},
						Usage:     "Create a new tag with a category and name",
						ArgsUsage: "<category> <name>",
						Action: func(c *cli.Context) error {
							category := c.Args().Get(0)
							name := c.Args().Get(1)
							ix.CreateTag(category, name, "./")
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					category := c.Args().Get(0)
					name := c.Args().Get(1)
					file := c.Args().Get(2)

					ix.Tag(category, name, file)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
