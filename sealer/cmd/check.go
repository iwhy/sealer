// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/check/service"

	"github.com/spf13/cobra"

	"github.com/alibaba/sealer/logger"
)

type CheckArgs struct {
	Pre  bool
	Post bool
}

var checkArgs *CheckArgs

// pushCmd represents the push command
var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "check the state of cluster ",
	Example: `sealer check --pre or sealer check --post`,
	Run: func(cmd *cobra.Command, args []string) {
		var checker service.CheckerService
		if checkArgs.Pre && checkArgs.Post {
			logger.Error("don't allow to set tow flags --pre and --post")
			os.Exit(1)
		} else if checkArgs.Pre {
			checker = service.NewPreCheckerService()
		} else if checkArgs.Post {
			checker = service.NewPostCheckerService()
		} else {
			checker = service.NewDefaultCheckerService()
		}
		if err := checker.Run(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	checkArgs = &CheckArgs{}
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&checkArgs.Pre, "pre", false, "Check dependencies before cluster creation")
	checkCmd.Flags().BoolVar(&checkArgs.Post, "post", false, "Check the status of the cluster after it is created")
}
