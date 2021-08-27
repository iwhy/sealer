// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/common"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	imageID   = "IMAGE ID"
	imageName = "IMAGE NAME"
)

var listCmd = &cobra.Command{
	Use:     "images",
	Short:   "list all cluster images",
	Example: `sealer images`,
	Run: func(cmd *cobra.Command, args []string) {
		ims, err := image.NewImageMetadataService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		imageMetadataList, err := ims.List()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		table := tablewriter.NewWriter(common.StdOut)
		table.SetHeader([]string{imageID, imageName})
		for _, image := range imageMetadataList {
			table.Append([]string{image.ID, image.Name})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
