// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/alibaba/sealer/logger"

	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/ssh"
)

var clusterfile string

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec commands on all hosts",
	Long:  `seautil exec -f Clusterfile "clean.sh"`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster := &v1.Cluster{}
		err := utils.UnmarshalYamlFile(clusterfile, cluster)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
		hosts := append(cluster.Spec.Masters.IPList, cluster.Spec.Nodes.IPList...)
		SSH := ssh.NewSSHByCluster(cluster)
		var wg sync.WaitGroup
		for _, host := range hosts {
			wg.Add(1)
			func(host string) {
				defer wg.Done()
				err := SSH.CmdAsync(host, args[0])
				if err != nil {
					logger.Error(err)
				}
			}(host)

		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&clusterfile, "clusterfile", "f", "", "cluster file filepath")
}
