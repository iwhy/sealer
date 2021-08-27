// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

var ImageTar string

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "save image",
	Long:  `save image to a file `,
	Example: `
sealer save -o [output file name] [image name]
save kubernetes:v1.18.3 image to kubernetes.tar.gz file:
sealer save -o kubernetes.tar.gz kubernetes:v1.18.3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ifs, err := image.NewImageFileService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if err = ifs.Save(args[0], ImageTar); err != nil {
			logger.Error("failed to save image %s, err: %v", args[0], err)
			os.Exit(1)
		}
		logger.Info("save image %s to %s successfully", args[0], ImageTar)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringVarP(&ImageTar, "output", "o", "", "write the image to a file")
	if err := saveCmd.MarkFlagRequired("output"); err != nil {
		logger.Error("failed to init flag: %v", err)
		os.Exit(1)
	}
}
