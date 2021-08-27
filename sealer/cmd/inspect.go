// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"fmt"
	"os"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

var clusterFilePrint bool

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "print the image information or clusterFile",
	Long: `sealer inspect kubernetes:v1.18.3 to print image information
sealer inspect -c kubernetes:v1.18.3 to print image Clusterfile`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if clusterFilePrint {
			cluster := image.GetClusterFileFromImageManifest(args[0])
			if cluster == "" {
				logger.Error("failed to find Clusterfile by image %s", args[0])
				os.Exit(1)
			}
			fmt.Println(cluster)
		} else {
			file, err := image.GetYamlByImage(args[0])
			if err != nil {
				logger.Error("failed to find information by image %s", args[0])
				os.Exit(1)
			}
			fmt.Println(file)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&clusterFilePrint, "Clusterfile", "c", false, "print the clusterFile")
}
