// Copyright © 2021 Alibaba Group Holding Ltd.

package apply

import (
	"fmt"

	"github.com/alibaba/sealer/common"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
)

type Interface interface {
	Apply() error
	Delete() error
}

func NewApplierFromFile(clusterfile string) (Interface, error) {
	cluster, err := GetClusterFromFile(clusterfile)
	if err != nil {
		return nil, err
	}
	return NewApplier(cluster)
}

func GetClusterFromFile(filepath string) (cluster *v1.Cluster, err error) {
	cluster = &v1.Cluster{}
	if err = utils.UnmarshalYamlFile(filepath, cluster); err != nil {
		return nil, fmt.Errorf("failed to get cluster from %s, %v", filepath, err)
	}
	cluster.SetAnnotations(common.ClusterfileName, filepath)
	return cluster, nil
}

func NewApplier(cluster *v1.Cluster) (Interface, error) {
	switch cluster.Spec.Provider {
	case common.AliCloud:
		return NewAliCloudProvider(cluster)
	case common.CONTAINER:
		return NewAliCloudProvider(cluster)
	}

	return NewDefaultApplier(cluster)
}
