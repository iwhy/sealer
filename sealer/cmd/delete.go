// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/alibaba/sealer/apply"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete a cluster",
	Long:    `if provider is BARESERVER will delete kubernetes nodes, or if provider is ALI_CLOUD, will delete all the infra resources`,
	Example: `sealer delete -f /root/.sealer/mycluster/Clusterfile [--force]`,
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if !force {
			var yesRx = regexp.MustCompile("^(?:y(?:es)?)$")
			var noRx = regexp.MustCompile("^(?:n(?:o)?)$")
			var input string
			for {
				fmt.Printf("Are you sure to delete the cluster? Yes [y/yes], No [n/no] : ")
				fmt.Scanln(&input)
				if yesRx.MatchString(input) {
					break
				}
				if noRx.MatchString(input) {
					fmt.Println("You have canceled to delete the cluster!")
					os.Exit(0)
				}
			}
		}
		applier, err := apply.NewApplierFromFile(clusterFile)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		if err = applier.Delete(); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&clusterFile, "Clusterfile", "f", "Clusterfile", "delete a kubernetes cluster with Clusterfile Annotations")
	deleteCmd.Flags().BoolP("force", "", false, "We also can input an --force flag to delete cluster by force")
}
