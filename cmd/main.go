package main

import (
	"fmt"
	"os"

	"github.com/SteveCastle/ix"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ix"
	app.Usage = "IX is a simple file indexer that uses a simple file structure to store metadata about files."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "store",
			Aliases: []string{"s"},
			Usage:   "Store returns the store location that will be used for saving metadata.",
			Action: func(c *cli.Context) error {
				fmt.Println(ix.FindStore("./"))
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:    "init",
					Aliases: []string{"i"},
					Usage:   "Init creates a new index in the current directory. Commands that run in directories under this index will be able to access the index.",
					Action: func(c *cli.Context) error {
						ix.InitIndex()
						return nil
					},
				},
			},
		},
		{
			Name:      "tag",
			Aliases:   []string{"t"},
			Usage:     "Tag a file with a category and name. If a directory is specified, all files in the directory will be tagged.",
			ArgsUsage: "<category> <name> <file>",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Aliases:   []string{"c"},
					Usage:     "Create a new tag with a category and name",
					ArgsUsage: "<category> <name>",
					Action: func(c *cli.Context) error {
						category := c.Args().Get(0)
						name := c.Args().Get(1)
						ix.CreateTag(category, name)
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
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	app.Run(os.Args)
}
