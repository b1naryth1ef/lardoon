package main

import (
	"log"
	"os"

	"github.com/b1naryth1ef/lardoon"
	"github.com/urfave/cli/v2"
)

func doServe(c *cli.Context) error {
	err := lardoon.InitDatabase(c.Path("db"))
	if err != nil {
		return err
	}

	var server lardoon.HTTPServer
	return server.Run(c.String("bind"))
}

func doImport(c *cli.Context) error {
	err := lardoon.InitDatabase(c.Path("db"))
	if err != nil {
		return err
	}

	return lardoon.ImportPath(c.Path("import-path"))
}

func doPrune(c *cli.Context) error {
	err := lardoon.InitDatabase(c.Path("db"))
	if err != nil {
		return err
	}

	return lardoon.PruneReplays(!c.Bool("no-dry-run"))
}

func main() {
	app := &cli.App{
		Name:        "lardoon",
		Description: "tacview repository",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "db",
				Usage:   "path to sqlite3 database file",
				Value:   "lardoon.db",
				EnvVars: []string{"LARDOON_DB_PATH"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "prune",
				Action: doPrune,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-dry-run",
						Usage: "during a dry-run no data will be mutated",
						Value: false,
					},
				},
			},
			{
				Name:   "import",
				Action: doImport,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "import-path",
						Usage:    "directory or replay path to import",
						Required: true,
						Aliases:  []string{"p"},
					},
				},
			},
			{
				Name:   "serve",
				Action: doServe,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "bind",
						Usage: "hostname/port to bind the server on",
						Value: "localhost:3883",
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
