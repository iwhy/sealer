// Copyright © 2021 Alibaba Group Holding Ltd.

package cache

import (
	"fmt"

	"github.com/alibaba/sealer/logger"

	"github.com/opencontainers/go-digest"

	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/pkg/errors"
)

//LocalImageCache saves all the layer
type LocalImageCache struct {
	chainStore ChainStore
}

type NopImageCache struct {
}

func NewLocalImageCache(chainStore ChainStore) (*LocalImageCache, error) {
	return &LocalImageCache{
		chainStore: chainStore,
	}, nil
}

func (NopImageCache) GetCache(parentID string, layer *Layer) (layerID digest.Digest, err error) {
	return "", errors.Errorf("nop cache")
}

func (lic *LocalImageCache) GetCache(parentID string, layer *Layer) (layerID digest.Digest, err error) {
	curChainID, err := layer.ChainID(ChainID(parentID))
	if err != nil {
		return "", fmt.Errorf("failed to get cur chain id, err: %s", err)
	}
	logger.Debug("current layer %+v, chain id %s", layer, curChainID)

	tmpLayer, err := getLocalCachedImage(lic.chainStore, curChainID)
	if err != nil {
		return "", err
	}

	return tmpLayer.ID, nil
}

func getLocalCachedImage(imageChain ChainStore, layerChainID ChainID) (v1.Layer, error) {
	return imageChain.GetChainLayer(layerChainID)
}
