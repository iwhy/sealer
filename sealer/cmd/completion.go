// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate autocompletion script for bash",
	Long: `Generate the autocompletion script for sealer for the bash shell.
To load completions in your current shell session:

	source <(sealer completion bash)

To load completions for every new session, execute once:

- Linux :
	## If bash-completion is not installed on Linux, please install the 'bash-completion' package
		sealer completion bash > /etc/bash_completion.d/sealer
	`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(common.StdOut); err != nil {
				logger.Error("failed to use bash completion, %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
