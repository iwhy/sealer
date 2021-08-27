// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import (
	"fmt"

	"github.com/alibaba/sealer/image/cache"
	"github.com/alibaba/sealer/logger"
	"github.com/opencontainers/go-digest"
)

type Prober interface {
	Reset()
	Probe(parentID string, layer *cache.Layer) (cacheID digest.Digest, err error)
}

type imageProber struct {
	cache       Cache
	reset       func() Cache
	cacheBusted bool
}

func NewImageProber(cacheBuilder CacheBuilder, noCache bool) Prober {
	if noCache {
		return &nopProber{}
	}

	reset := func() Cache {
		c, err := cacheBuilder.BuildImageCache()
		if err != nil {
			logger.Info("failed to init image cache, err: %s", err)
			return &cache.NopImageCache{}
		}
		return c
	}

	return &imageProber{cache: reset(), reset: reset}
}

func (c *imageProber) Reset() {
	c.cache = c.reset()
	c.cacheBusted = false
}

func (c *imageProber) Probe(parentID string, layer *cache.Layer) (cacheID digest.Digest, err error) {
	if c.cacheBusted {
		return "", nil
	}

	cacheID, err = c.cache.GetCache(parentID, layer)
	if err != nil {
		return "", err
	}

	return cacheID, nil
}

type nopProber struct{}

func (c *nopProber) Reset() {}

func (c *nopProber) Probe(_ string, _ *cache.Layer) (digest.Digest, error) {
	return "", fmt.Errorf("nop prober")
}
