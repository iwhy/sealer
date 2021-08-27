// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:     "tag",
	Short:   "tag IMAGE[:TAG] TARGET_IMAGE[:TAG]",
	Example: `sealer tag sealer/cloudrootfs:v1.16.9-alpha.6 registry.cn-qingdao.aliyuncs.com/sealer-io/cloudrootfs:v1.16.9-alpha.5`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ims, err := image.NewImageMetadataService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		err = ims.Tag(args[0], args[1])
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
