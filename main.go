package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "acdfuse"
	app.Usage = "mount amazon cloud drive as fuse file system"
	app.Action = func(c *cli.Context) {
		println("boom! I say!")
	}

	app.Run(os.Args)
}
