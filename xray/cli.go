package xray

import (
	"github.com/codegangsta/cli"
	"github.com/jfrogdev/jfrog-cli-go/xray/commands"
	"github.com/jfrogdev/jfrog-cli-go/utils/cliutils"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "offline-update",
			Usage:   "Download Xray offline updates",
			Flags:   offlineUpdateFlags(),
			Action: offlineUpdates,
		},
	}
}

func offlineUpdateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "[Mandatory] Global component database url",
		},
		cli.StringFlag{
			Name:  "license-id",
			Usage: "[Mandatory] Xray license ID",
		},
	}
}

func getOfflineUpdatesFlag(c *cli.Context) (flags *commands.OfflineUpdatesFlags) {
	flags = new(commands.OfflineUpdatesFlags)
	flags.Url = c.String("url");
	flags.License = c.String("license-id");
	if len(flags.License) < 1 || len(flags.Url) < 1 {
		cliutils.Exit(cliutils.ExitCodeError, "Url and license-id are mandatory arguments.")
	}
	return
}

func offlineUpdates(c *cli.Context) {
	offlineUpdateFlags := getOfflineUpdatesFlag(c)
	err := commands.OfflineUpdate(offlineUpdateFlags)
	if err != nil {
		cliutils.Exit(cliutils.ExitCodeError, err.Error())
	}
}
