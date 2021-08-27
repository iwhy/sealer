// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/version"
	"github.com/spf13/cobra"
)

var shortPrint bool

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "version",
	Example: `sealer version`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalled, err := json.Marshal(version.Get())
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if shortPrint {
			fmt.Println(version.Get().String())
		} else {
			fmt.Println(string(marshalled))
		}

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&shortPrint, "short", false, "If true, print just the version number.")
}
