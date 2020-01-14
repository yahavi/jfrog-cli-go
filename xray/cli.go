package xray

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-go/docs/common"
	"github.com/jfrog/jfrog-cli-go/docs/xray/offlineupdate"
	"github.com/jfrog/jfrog-cli-go/xray/commands"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"time"
)

const DATE_FORMAT = "2006-01-02"

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:         "offline-update",
			Usage:        offlineupdate.Description,
			HelpName:     common.CreateUsage("xr offline-update", offlineupdate.Description, offlineupdate.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			Flags:        offlineUpdateFlags(),
			Aliases:      []string{"ou"},
			BashComplete: common.CreateBashCompletionFunc(),
			Action:       offlineUpdates,
		},
	}
}

func offlineUpdateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "license-id",
			Usage: "[Mandatory] Xray license ID` `",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "[Optional] From update date in YYYY-MM-DD format.` `",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "[Optional] To update date in YYYY-MM-DD format.` `",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "[Optional] Xray API version.` `",
		},
	}
}

func getOfflineUpdatesFlag(c *cli.Context) (flags *commands.OfflineUpdatesFlags, err error) {
	flags = new(commands.OfflineUpdatesFlags)
	flags.Version = c.String("version")
	flags.License = c.String("license-id")
	if len(flags.License) < 1 {
		return nil, errors.New("The --license-id option is mandatory.")
	}
	from := c.String("from")
	to := c.String("to")
	if len(to) > 0 && len(from) < 1 {
		return nil, errors.New("The --from option is mandatory, when the --to option is sent.")
	}
	if len(from) > 0 && len(to) < 1 {
		return nil, errors.New("The --to option is mandatory, when the --from option is sent.")
	}
	if len(from) > 0 && len(to) > 0 {
		flags.From, err = dateToMilliseconds(from)
		errorutils.WrapError(err)
		if err != nil {
			return
		}
		flags.To, err = dateToMilliseconds(to)
		errorutils.WrapError(err)
	}
	return
}

func dateToMilliseconds(date string) (dateInMillisecond int64, err error) {
	t, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		errorutils.WrapError(err)
		return
	}
	dateInMillisecond = t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	return
}

func offlineUpdates(c *cli.Context) error {
	offlineUpdateFlags, err := getOfflineUpdatesFlag(c)
	if err != nil {
		return err
	}

	return commands.OfflineUpdate(offlineUpdateFlags)
}
