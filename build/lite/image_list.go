// Copyright © 2021 Alibaba Group Holding Ltd.

package lite

type Interface interface {
	// List all the containers images in helm charts
	ListImages(clusterName string) ([]string, error)
}
