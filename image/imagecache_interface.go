// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import (
	"github.com/alibaba/sealer/image/cache"
	"github.com/opencontainers/go-digest"
)

type CacheBuilder interface {
	BuildImageCache() (Cache, error)
}

type Cache interface {
	GetCache(parentID string, layer *cache.Layer) (LayerID digest.Digest, err error)
}
