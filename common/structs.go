// Copyright © 2021 Alibaba Group Holding Ltd.

package common

type RunArgs struct {
	Masters    string
	Nodes      string
	User       string
	Password   string
	Pk         string
	PkPassword string
	PodCidr    string
	SvcCidr    string
}
