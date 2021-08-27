// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/alibaba/sealer/apply"
	"github.com/alibaba/sealer/logger"
)

var clusterFile string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply a kubernetes cluster",
	Long:  `seautil apply -f cluster.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		applier, err := apply.NewApplierFromFile(clusterFile)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if err = applier.Apply(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&clusterFile, "clusterfile", "f", "", "cluster file filepath")
}
