package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jfrog/jfrog-cli-go/artifactory"
	"github.com/jfrog/jfrog-cli-go/bintray"
	"github.com/jfrog/jfrog-cli-go/completion"
	"github.com/jfrog/jfrog-cli-go/docs/common"
	"github.com/jfrog/jfrog-cli-go/missioncontrol"
	"github.com/jfrog/jfrog-cli-go/utils/cliutils"
	"github.com/jfrog/jfrog-cli-go/utils/log"
	"github.com/jfrog/jfrog-cli-go/xray"
	"github.com/jfrog/jfrog-client-go/utils"
)

const commandHelpTemplate string = `{{.HelpName}}{{if .UsageText}}
Arguments:
{{.UsageText}}
{{end}}{{if .Flags}}
Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}{{if .ArgsUsage}}
Environment Variables:
{{.ArgsUsage}}{{end}}

`

const appHelpTemplate string = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} [arguments...]{{end}}
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
Environment Variables:
` + common.GlobalEnvVars + `{{end}}

`

const subcommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.HelpName}} command{{if .Flags}} [command options]{{end}}[arguments...]

COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
Environment Variables:
` + common.GlobalEnvVars + `{{end}}

`

func main() {
	log.SetDefaultLogger()
	err := execMain()
	cliutils.ExitOnErr(err)
}

func execMain() error {
	// Set JFrog CLI's user-agent on the jfrog-client-go.
	utils.SetUserAgent(cliutils.GetUserAgent())

	app := cli.NewApp()
	app.Name = "jfrog"
	app.Usage = "See https://github.com/jfrog/jfrog-cli-go for usage instructions."
	app.Version = cliutils.GetVersion()
	args := os.Args
	app.EnableBashCompletion = true
	app.Commands = getCommands()
	app.Flags = getFlags()
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.AppHelpTemplate = appHelpTemplate
	cli.SubcommandHelpTemplate = subcommandHelpTemplate
	err := app.Run(args)
	return err
}

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        cliutils.CmdArtifactory,
			Usage:       "Artifactory commands",
			Subcommands: artifactory.GetCommands(),
		},
		{
			Name:        cliutils.CmdBintray,
			Usage:       "Bintray commands",
			Subcommands: bintray.GetCommands(),
		},
		{
			Name:        cliutils.CmdMissionControl,
			Usage:       "Mission Control commands",
			Subcommands: missioncontrol.GetCommands(),
		},
		{
			Name:        cliutils.CmdXray,
			Usage:       "Xray commands",
			Subcommands: xray.GetCommands(),
		},
		{
			Name:        cliutils.CmdCompletion,
			Usage:       "Generate autocomplete scripts",
			Subcommands: completion.GetCommands(),
		},
	}
}

func getFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "stacktrace",
			Usage: "[Default: false] Set to true to print the error stacktrace for all errors.` `",
		},
	}
}
