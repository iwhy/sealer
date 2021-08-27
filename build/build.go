// Copyright © 2021 Alibaba Group Holding Ltd.

package build

import "github.com/alibaba/sealer/common"

type Interface interface {
	Build(name string, context string, kubefileName string) error
}

func NewBuilder(config *Config) (Interface, error) {
	switch config.BuildType {
	case common.LiteBuild:
		return NewLiteBuilder(config)
	case common.LocalBuild:
		return NewLocalBuilder(config)
	default:
		return NewCloudBuilder(config)
	}
}
