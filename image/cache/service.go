// Copyright © 2021 Alibaba Group Holding Ltd.

package cache

import (
	"fmt"

	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/opencontainers/go-digest"
)

type Service interface {
	NewCacheLayer(layer v1.Layer, cacheID digest.Digest) Layer

	CalculateChainID(layers interface{}) (ChainID, error)
}

type service struct {
}

func (s *service) NewCacheLayer(layer v1.Layer, cacheID digest.Digest) Layer {
	return Layer{
		CacheID: cacheID.String(),
		Type:    layer.Type,
		Value:   layer.Value,
	}
}

func (s *service) CalculateChainID(layers interface{}) (ChainID, error) {
	switch ls := layers.(type) {
	case []Layer:
		return CalculateCacheID(ls)
	default:
		return "", fmt.Errorf("do not support calculate chain ID on %v", ls)
	}
}

func NewService() (Service, error) {
	return &service{}, nil
}
