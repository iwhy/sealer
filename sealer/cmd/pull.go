// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"

	"os"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:     "pull",
	Short:   "pull cloud image to local",
	Example: `sealer pull registry.cn-qingdao.aliyuncs.com/sealer-io/cloudrootfs:v1.16.9-alpha.5`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imgSvc, err := image.NewImageService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		if err := imgSvc.Pull(args[0]); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		logger.Info("Pull %s success", args[0])
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
