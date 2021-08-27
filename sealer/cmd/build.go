// Copyright © 2021 Alibaba Group Holding Ltd.

package cmd

import (
	"github.com/alibaba/sealer/build"
	"github.com/alibaba/sealer/logger"
	"github.com/spf13/cobra"

	"os"
)

type BuildFlag struct {
	ImageName    string
	KubefileName string
	BuildType    string
	NoCache      bool
	Lite         bool
	ImageList    string
}

var buildConfig *BuildFlag

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build [flags] PATH",
	Short:   "cloud image local build command line",
	Example: `sealer build -f Kubefile -t my-kubernetes:1.18.3 .`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error("sealer build requires exactly 1 argument.")
			os.Exit(1)
		}
		conf := &build.Config{
			BuildType: buildConfig.BuildType,
			NoCache:   buildConfig.NoCache,
			ImageName: buildConfig.ImageName,
		}
		builder, err := build.NewBuilder(conf)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		err = builder.Build(buildConfig.ImageName, args[0], buildConfig.KubefileName)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	buildConfig = &BuildFlag{}
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&buildConfig.BuildType, "buildType", "b", "", "cluster image build type,default is cloud")
	buildCmd.Flags().StringVarP(&buildConfig.KubefileName, "kubefile", "f", "Kubefile", "kubefile filepath")
	buildCmd.Flags().StringVarP(&buildConfig.ImageName, "imageName", "t", "", "cluster image name")
	buildCmd.Flags().BoolVar(&buildConfig.NoCache, "no-cache", false, "build without cache")
	if err := buildCmd.MarkFlagRequired("imageName"); err != nil {
		logger.Error("failed to init flag: %v", err)
		os.Exit(1)
	}
}
