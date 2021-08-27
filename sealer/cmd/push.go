// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/image/utils"
	"github.com/alibaba/sealer/logger"

	"os"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:     "push",
	Short:   "push cloud image to registry",
	Example: `sealer push registry.cn-qingdao.aliyuncs.com/sealer-io/my-kuberentes-cluster-with-dashboard:latest`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imgsvc, err := image.NewImageService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		if err = imgsvc.Push(args[0]); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
	ValidArgsFunction: utils.ImageListFuncForCompletion,
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
