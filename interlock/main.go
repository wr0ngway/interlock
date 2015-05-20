package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ehazlett/interlock/interlock/commands"
	"github.com/ehazlett/interlock/plugins"
	"github.com/ehazlett/interlock/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "interlock"
	app.Version = version.FULL_VERSION
	app.Author = "@ehazlett"
	app.Email = ""
	app.Usage = "event driven docker plugins"
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "swarm-url, s",
			Value: "unix:///var/run/docker.sock",
			Usage: "swarm addr",
		},
		cli.StringFlag{
			Name:  "swarm-tls-ca-cert",
			Value: "",
			Usage: "tls ca certificate",
		},
		cli.StringFlag{
			Name:  "swarm-tls-cert",
			Value: "",
			Usage: "tls certificate",
		},
		cli.StringFlag{
			Name:  "swarm-tls-key",
			Value: "",
			Usage: "tls key",
		},
		cli.BoolFlag{
			Name:  "swarm-allow-insecure",
			Usage: "enable insecure tls communication",
		},
		cli.StringSliceFlag{
			Name:  "plugin, p",
			Usage: "enable plugin",
			Value: &cli.StringSlice{},
		},
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug",
		},
	}
	// base commands
	baseCommands := []cli.Command{
		commands.CmdStart,
		commands.CmdListPlugins,
	}
	// plugin supplied commands
	baseCommands = append(baseCommands, plugins.GetCommands()...)

	app.Commands = baseCommands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
