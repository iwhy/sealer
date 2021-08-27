// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"github.com/sealyun/lvscare/care"
	"github.com/spf13/cobra"
)

var Ipvs care.LvsCare

// ipvsCmd represents the ipvs command
var ipvsCmd = &cobra.Command{
	Use:   "ipvs",
	Short: "seautil create or care local ipvs LB",
	Long: `create ipvs rules: seautil ipvs --vs 10.1.1.2:6443 --rs 192.168.0.2:6443 --rs 192.168.0.3:6443 --health-path /healthz --health-schem https --run-once
clean ipvs rules: seautil ipvs clean`,
	Run: func(cmd *cobra.Command, args []string) {
		Ipvs.VsAndRsCare()
	},
}

func init() {
	rootCmd.AddCommand(ipvsCmd)
	ipvsCmd.Flags().BoolVar(&Ipvs.RunOnce, "run-once", false, "run once mode")
	ipvsCmd.Flags().BoolVarP(&Ipvs.Clean, "clean", "c", true, " clean Vip ipvs rule before join node, if Vip has no ipvs rule do nothing.")
	ipvsCmd.Flags().StringVar(&Ipvs.VirtualServer, "vs", "", "virturl server like 10.54.0.2:6443")
	ipvsCmd.Flags().StringSliceVar(&Ipvs.RealServer, "rs", []string{}, "virturl server like 192.168.0.2:6443")

	ipvsCmd.Flags().StringVar(&Ipvs.HealthPath, "health-path", "/healthz", "health check path")
	ipvsCmd.Flags().StringVar(&Ipvs.HealthSchem, "health-schem", "https", "health check schem")
	ipvsCmd.Flags().Int32Var(&Ipvs.Interval, "interval", 5, "health check interval, unit is sec.")
}
