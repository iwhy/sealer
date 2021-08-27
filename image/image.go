// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import "github.com/alibaba/sealer/image/store"

func NewImageService() (Service, error) {
	imageStore, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}

	return DefaultImageService{imageStore: imageStore}, nil
}

func NewImageMetadataService() (MetadataService, error) {
	imageStore, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}
	return DefaultImageMetadataService{
		imageStore: imageStore,
	}, nil
}

func NewImageFileService() (FileService, error) {
	layerStore, err := store.NewDefaultLayerStore()
	if err != nil {
		return nil, err
	}

	imageStore, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}
	return DefaultImageFileService{
		layerStore: layerStore,
		imageStore: imageStore,
	}, nil
}

func NewDeleteImageService(force bool) (Service, error) {
	imageStore, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}
	return DefaultImageService{
		imageStore:       imageStore,
		ForceDeleteImage: force,
	}, nil
}
