// Copyright © 2021 Alibaba Group Holding Ltd.

package image

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/alibaba/sealer/image/store"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/mount"
)

// GetImageLayerDirs return image hash list
// current image is different with the image in build stage
// current image has no from layer
func GetImageLayerDirs(image *v1.Image) (res []string, err error) {
	for _, layer := range image.Spec.Layers {
		if layer.ID != "" {
			res = append(res, filepath.Join(common.DefaultLayerDir, layer.ID.Hex()))
		}
	}
	return
}

// GetClusterFileFromImage retrieve ClusterFile From image
func GetClusterFileFromImage(imageName string) string {
	clusterfile := GetClusterFileFromImageManifest(imageName)
	if clusterfile != "" {
		return clusterfile
	}

	clusterfile = GetFileFromBaseImage(imageName, "etc", common.DefaultClusterFileName)
	if clusterfile != "" {
		return clusterfile
	}
	return ""
}

// GetClusterFileFromImageManifest retrieve ClusterFile from image manifest(image yaml)
func GetClusterFileFromImageManifest(imageName string) string {
	//  find cluster file from image manifest
	var (
		image *v1.Image
		err   error
	)
	is, err := store.NewDefaultImageStore()
	if err != nil {
		logger.Error("failed to init image store, err: %s", err)
		return ""
	}
	image, err = is.GetByName(imageName)
	if err != nil {
		ims, err := NewImageMetadataService()
		if err != nil {
			logger.Error("failed to create image metadata svc, err: %v", err)
		}

		imageMetadata, err := ims.GetRemoteImage(imageName)
		if err != nil {
			logger.Error("failed to find image %s,err: %v", imageName, err)
			return ""
		}
		image = &imageMetadata
	}
	Clusterfile, ok := image.Annotations[common.ImageAnnotationForClusterfile]
	if !ok {
		logger.Error("failed to find Clusterfile in local")
		return ""
	}
	return Clusterfile
}

// GetFileFromBaseImage retrieve file from base image
func GetFileFromBaseImage(imageName string, paths ...string) string {
	mountTarget, _ := utils.MkTmpdir()
	mountUpper, _ := utils.MkTmpdir()
	defer func() {
		utils.CleanDirs(mountTarget, mountUpper)
	}()

	imgSvc, err := NewImageService()
	if err != nil {
		return ""
	}
	if err = imgSvc.PullIfNotExist(imageName); err != nil {
		return ""
	}

	driver := mount.NewMountDriver()
	is, err := store.NewDefaultImageStore()
	if err != nil {
		logger.Error("failed to init image store, err: %s", err)
		return ""
	}
	image, err := is.GetByName(imageName)
	if err != nil {
		return ""
	}

	layers, err := GetImageLayerDirs(image)
	if err != nil {
		return ""
	}

	err = driver.Mount(mountTarget, mountUpper, layers...)
	if err != nil {
		return ""
	}
	defer func() {
		err := driver.Unmount(mountTarget)
		if err != nil {
			logger.Warn(err)
		}
	}()
	var subPath []string
	subPath = append(subPath, mountTarget)
	subPath = append(subPath, paths...)
	clusterFile := filepath.Join(subPath...)
	data, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		return ""
	}
	return string(data)
}

func GetYamlByImage(imageName string) (string, error) {
	img, err := GetImageByName(imageName)
	if err != nil {
		return "", fmt.Errorf("failed to get image %s, err: %s", imageName, err)
	}

	ImageInformation, err := yaml.Marshal(img)
	if err != nil {
		return "", err
	}

	return string(ImageInformation), nil
}

func GetImageByName(imageName string) (*v1.Image, error) {
	is, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, fmt.Errorf("failed to init image store, err: %s", err)
	}
	img, err := is.GetByName(imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s, err: %s", imageName, err)
	}
	return img, nil
}
