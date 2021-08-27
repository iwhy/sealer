// Copyright © 2021 Alibaba Group Holding Ltd.

package distributionutil

import (
	"time"

	"github.com/alibaba/sealer/image/reference"
	"github.com/alibaba/sealer/image/store"
	"github.com/docker/docker/pkg/progress"
)

type Config struct {
	LayerStore     store.LayerStore
	ProgressOutput progress.Output
	Named          reference.Named
}

type registryConfig struct {
	Domain   string
	Insecure bool
	SkipPing bool
	NonSSL   bool
	Timeout  time.Duration
	Headers  map[string]string
}
