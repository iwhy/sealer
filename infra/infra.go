// Copyright © 2021 Alibaba Group Holding Ltd.

package infra

import (
	"fmt"

	"github.com/alibaba/sealer/infra/container"

	"github.com/alibaba/sealer/infra/aliyun"
	v1 "github.com/alibaba/sealer/types/api/v1"
)

type Interface interface {
	// Apply apply iaas resources and save metadata info like vpc instance id to cluster status
	// https://github.com/fanux/sealgate/tree/master/cloud
	Apply() error
}

func NewDefaultProvider(cluster *v1.Cluster) (Interface, error) {
	switch cluster.Spec.Provider {
	case aliyun.AliCloud:
		return NewAliProvider(cluster)
	case container.CONTAINER:
		return NewContainerProvider(cluster)
	default:
		return nil, fmt.Errorf("the provider is invalid")
	}
}

func NewAliProvider(cluster *v1.Cluster) (Interface, error) {
	config := new(aliyun.Config)
	err := aliyun.LoadConfig(config)
	if err != nil {
		return nil, err
	}
	aliProvider := new(aliyun.AliProvider)
	aliProvider.Config = *config
	aliProvider.Cluster = cluster
	err = aliProvider.NewClient()
	if err != nil {
		return nil, err
	}
	return aliProvider, nil
}

func NewContainerProvider(cluster *v1.Cluster) (Interface, error) {
	if container.IsDockerAvailable() {
		return nil, fmt.Errorf("please install docker on your system")
	}

	cli, err := container.NewClientWithCluster(cluster)
	if err != nil {
		return nil, fmt.Errorf("new container client failed")
	}

	return cli, nil
}
