// Copyright © 2021 Alibaba Group Holding Ltd.

package store

import (
	"io"

	"github.com/alibaba/sealer/image/reference"

	"github.com/opencontainers/go-digest"
)

type LayerStore interface {
	Get(id LayerID) Layer
	RegisterLayerIfNotPresent(layer Layer) error
	RegisterLayerForBuilder(diffPath string) (digest.Digest, error)
	Delete(id LayerID) error
	DisassembleTar(layerID digest.Digest, streamReader io.ReadCloser) error
	AddDistributionMetadata(layerID LayerID, named reference.Named, descriptorDigest digest.Digest) error
}

type Layer interface {
	ID() LayerID
	TarStream() (io.ReadCloser, error)
	SimpleID() string
	Size() int64
	MediaType() string
	DistributionMetadata() map[string]digest.Digest
	SetSize(size int64)
}
