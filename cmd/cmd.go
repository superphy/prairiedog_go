package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/superphy/prairiedog/graph"
)

var Version = "0.1"
var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}
DESCRIPTION:
{{.Description}}
USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]
COMMANDS:
{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{if .Flags}}
FLAGS:
{{range .Flags}}{{.}}
{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "i",
		Usage: "path to input files",
		Value: "",
	},
	cli.StringFlag{
		Name:  "temp-path",
		Usage: "path to temp files",
		Value: os.TempDir(),
	},
	cli.StringFlag{
		Name:  "log",
		Usage: "/var/log/prairiedog.log",
		Value: "/var/log/prairiedog.log",
	},
	cli.StringFlag{
		Name:  "port",
		Usage: "port for backing database",
		Value: "8000",
	},
	cli.StringFlag{
		Name:  "address",
		Usage: "localhost",
		Value: "localhost",
	},
}

type Cmd struct {
	*cli.App
}

func VersionAction(c *cli.Context) {
	fmt.Println(color.YellowString(fmt.Sprintf("prairiedog: Pangenome graphs")))
}

func New() *Cmd {
	app := cli.NewApp()
	app.Name = "prairiedog"
	app.Author = ""
	app.Usage = "prairiedog"
	app.Description = `Pangenome graph`
	app.Flags = globalFlags
	app.CustomAppHelpTemplate = helpTemplate
	app.Commands = []cli.Command{
		{
			Name:   "version",
			Action: VersionAction,
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) {
		options := []graph.OptionFn{}
		if v := c.String("i"); v != "" {
			options = append(options, graph.InputFiles(v))
		}
		if v := c.String("address"); v != "" {
			options = append(options, graph.Address(s))
		}
		if v := c.String("port"); v != "" {
			options = append(options, graph.Port(s))
		}
		if v := c.String("logfile"); v != "" {
			options = append(options, graph.LogFile(s))
		}

		g, err := graph.New(
			options...,
		)

		if err != nil {
			fmt.Println(color.RedString("Error running prairiedog: %s", err.Error()))
			return
		}

		g.Run()
	}

	return &Cmd{
		App: app,
	}
}
