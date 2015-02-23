package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/homedir"
)

var (
	defaultPluginPath = filepath.Join(homedir.Get(), ".interlock")
)

var appCommands = []cli.Command{
	{
		Name:   "ls",
		Usage:  "list available plugins",
		Action: listPlugins,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "plugin-path, p",
				Usage:  "path for plugins",
				Value:  defaultPluginPath,
				EnvVar: "INTERLOCK_PLUGIN_PATH",
			},
		},
	},
}

func listPlugins(c *cli.Context) {
	pluginPath := c.String("plugin-path")

	plugins, err := ioutil.ReadDir(pluginPath)
	if err != nil {

	}

	for _, fi := range plugins {
		println(fi.Name())
	}
}
