package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/ehazlett/interlock/plugins"
)

var CmdListPlugins = cli.Command{
	Name:   "list-plugins",
	Usage:  "List available plugins",
	Action: cmdListPlugins,
}

func cmdListPlugins(c *cli.Context) {
	allPlugins := plugins.GetPlugins()
	w := tabwriter.NewWriter(os.Stdout, 8, 1, 3, ' ', 0)

	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION\tURL")

	for _, p := range allPlugins {
		i := p.Info()
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			i.Name,
			i.Version,
			i.Description,
			i.Url,
		)
	}
	w.Flush()
}
