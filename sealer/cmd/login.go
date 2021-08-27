// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"os"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"
)

type LoginFlag struct {
	RegistryURL      string
	RegistryUsername string
	RegistryPasswd   string
}

var loginConfig *LoginFlag

var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "login image repositories",
	Example: `sealer login registry.cn-qingdao.aliyuncs.com -u [username] -p [password]`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error("enter the registry URL failed")
			os.Exit(1)
		}

		imgSvc, err := image.NewImageService()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		if err = imgSvc.Login(args[0], loginConfig.RegistryUsername, loginConfig.RegistryPasswd); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	loginConfig = &LoginFlag{}
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginConfig.RegistryUsername, "username", "u", "", "user name for login registry")
	loginCmd.Flags().StringVarP(&loginConfig.RegistryPasswd, "passwd", "p", "", "password for login registry")
	if err := loginCmd.MarkFlagRequired("username"); err != nil {
		logger.Error("failed to init flag: %v", err)
		os.Exit(1)
	}
	if err := loginCmd.MarkFlagRequired("passwd"); err != nil {
		logger.Error("failed to init flag: %v", err)
		os.Exit(1)
	}
}
