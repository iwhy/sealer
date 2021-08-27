// Copyright © 2021 Alibaba Group Holding Ltd.

package cache

import (
	"strings"

	"github.com/opencontainers/go-digest"
)

type Layer struct {
	// cacheID for layerdb/layersid/cacheID, we will load the content only for COPY layer
	CacheID string `json:"cache_id"`
	// same as v1Layer type
	Type string `json:"type"`
	// same as v1Layer value
	Value string `json:"value"`
}

func (l *Layer) String() string {
	return strings.TrimSpace(l.CacheID) + ":" + strings.TrimSpace(l.Type) + ":" + strings.TrimSpace(l.Value)
}

func (l *Layer) ChainID(parentID ChainID) (ChainID, error) {
	if parentID.String() == "" {
		return ChainID(digest.FromString(l.String())), nil
	}
	return ChainID(digest.FromString(parentID.String() + ":" + l.String())), nil
}
