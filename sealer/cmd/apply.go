// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"github.com/alibaba/sealer/apply"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"

	"os"
)

var clusterFile string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:     "apply",
	Short:   "apply a kubernetes cluster",
	Example: `sealer apply -f Clusterfile`,
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
	applyCmd.Flags().StringVarP(&clusterFile, "Clusterfile", "f", "Clusterfile", "apply a kubernetes cluster")
}
