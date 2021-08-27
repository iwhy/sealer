// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import (
	"fmt"

	"github.com/alibaba/sealer/image/cache"
	"github.com/alibaba/sealer/image/store"
)

func (d DefaultImageService) BuildImageCache() (Cache, error) {
	ls, err := store.NewDefaultLayerStore()
	if err != nil {
		return nil, fmt.Errorf("failed to build image cache, err: %s", err)
	}
	fs, err := store.NewFSStoreBackend()
	if err != nil {
		return nil, fmt.Errorf("failed to init store backend for image cache, err: %s", err)
	}
	imageStore, err := cache.NewImageStore(fs, ls)
	if err != nil {
		return nil, fmt.Errorf("failed to init image store for image cache, err: %s", err)
	}

	return cache.NewLocalImageCache(imageStore)
}
