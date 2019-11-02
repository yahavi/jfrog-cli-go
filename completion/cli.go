package completion

import (
	"github.com/codegangsta/cli"
	bash "github.com/jfrog/jfrog-cli-go/completion/bash"
	zsh "github.com/jfrog/jfrog-cli-go/completion/zsh"
	"github.com/jfrog/jfrog-cli-go/docs/common"
	bashdoc "github.com/jfrog/jfrog-cli-go/docs/completion/bash"
	zshdoc "github.com/jfrog/jfrog-cli-go/docs/completion/zsh"
)

func GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:         "bash",
			Usage:        bashdoc.Description,
			HelpName:     common.CreateUsage("completion bash", bashdoc.Description, bashdoc.Usage),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(*cli.Context) {
				bash.WriteBashCompletionScript()
			},
		},
		{
			Name:         "zsh",
			Usage:        zshdoc.Description,
			HelpName:     common.CreateUsage("completion zsh", zshdoc.Description, zshdoc.Usage),
			BashComplete: common.CreateBashCompletionFunc(),
			Action: func(*cli.Context) {
				zsh.WriteZshCompletionScript()
			},
		},
	}
}
