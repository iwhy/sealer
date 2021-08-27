// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import (
	"github.com/alibaba/sealer/image/store"
	"github.com/alibaba/sealer/image/types"
	v1 "github.com/alibaba/sealer/types/api/v1"
)

// MetadataService is the interface for providing image metadata service
type MetadataService interface {
	Tag(imageName, tarImageName string) error
	List() ([]types.ImageMetadata, error)
	GetImage(imageName string) (*v1.Image, error)
	GetRemoteImage(imageName string) (v1.Image, error)
	DeleteImage(imageName string) error
}

// FileService is the interface for file operations
type FileService interface {
	Load(imageSrc string) error
	Save(imageName string, imageTar string) error
	Merge(image *v1.Image) error
}

// Service is image service
type Service interface {
	Pull(imageName string) error
	PullIfNotExist(imageName string) error
	Push(imageName string) error
	Delete(imageName string) error
	Login(RegistryURL, RegistryUsername, RegistryPasswd string) error
	CacheBuilder
}

type LayerService interface {
	LayerStore() store.LayerStore
}
