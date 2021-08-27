// Copyright © 2021 Alibaba Group Holding Ltd.

package distributionutil

import (
	"context"

	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/utils"

	"github.com/alibaba/sealer/image/reference"

	"github.com/docker/distribution"
)

func NewV2Repository(named reference.Named, actions ...string) (distribution.Repository, error) {
	authConfig, err := utils.GetDockerAuthInfoFromDocker(named.Domain())
	if err != nil {
		logger.Warn("failed to get auth info, err: %s", err)
	}

	repo, err := NewRepository(context.Background(), authConfig, named.Repo(), registryConfig{Insecure: true, Domain: named.Domain()}, actions...)
	if err == nil {
		return repo, nil
	}

	return NewRepository(context.Background(), authConfig, named.Repo(), registryConfig{Insecure: true, NonSSL: true, Domain: named.Domain()}, actions...)
}
