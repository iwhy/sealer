// Copyright © 2021 Alibaba Group Holding Ltd.

package build

import (
	"fmt"
	"path/filepath"

	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/utils/mount"

	manifest "github.com/alibaba/sealer/build/lite/manifests"

	"github.com/alibaba/sealer/build/lite/charts"
	"github.com/alibaba/sealer/build/lite/docker"
	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/filesystem"
	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
)

type LiteBuilder struct {
	local *LocalBuilder
}

func (me *LiteBuilder) Build(name string, context string, kubefileName string) error {
	err := me.local.initBuilder(name, context, kubefileName)
	if err != nil {
		return err
	}

	pipLine, err := me.GetBuildPipeLine()
	if err != nil {
		return err
	}

	for _, f := range pipLine {
		if err = f(); err != nil {
			return err
		}
	}
	return nil
}

//实际上，这里依然执行了local_builder的三个步骤
// 外加上， lite本身的8个步骤，总共11个步骤
func (me *LiteBuilder) GetBuildPipeLine() ([]func() error, error) {
	var buildPipeline []func() error
	if err := me.local.InitImageSpec(); err != nil {
		return nil, err
	}

	buildPipeline = append(buildPipeline,
		me.PreCheck,
		me.local.PullBaseImageNotExist,
		me.InitClusterFile,
		me.MountImage,
		me.local.ExecBuild,
		me.local.UpdateImageMetadata,
		me.ReMountImage,
		me.InitDockerAndRegistry,
		me.CacheImageToRegistry,
		me.AddUpperLayerToImage,
		me.Clear,
	)
	return buildPipeline, nil
}


//#region pipline 8个步骤

// step 1： 检查本地docker images是否已经存在同名镜像
// 镜像列表是通过"github.com/docker/docker/client"客户端获取的
// @todo: 因此这样就可以在代码里面，执行很多docker操作
func (me *LiteBuilder) PreCheck() error {
	d := docker.Docker{}
	images, _ := d.ImagesList()
	if len(images) > 0 {
		logger.Warn("The image already exists on the host. Note that the existing image cannot be cached in registry")
	}
	return nil
}

// step 3
// load cluster file from disk\
// 没有被引用，，
func (me *LiteBuilder) InitClusterFile() error {
	clusterFile := common.TmpClusterfile
	if !utils.IsFileExist(clusterFile) {
		rawClusterFile := GetRawClusterFile(me.local.Image)
		if rawClusterFile == "" {
			return fmt.Errorf("failed to get cluster file from context or base image")
		}
		err := utils.WriteFile(common.RawClusterfile, []byte(rawClusterFile))
		if err != nil {
			return err
		}
		clusterFile = common.RawClusterfile
	}

	var cluster v1.Cluster
	err := utils.UnmarshalYamlFile(clusterFile, &cluster)
	if err != nil {
		return fmt.Errorf("failed to read %s:%v", clusterFile, err)
	}
	me.local.Cluster = &cluster

	logger.Info("read cluster file %s success !", clusterFile)
	return nil
}

// step 4
func (me *LiteBuilder) MountImage() error {
	FileSystem, err := filesystem.NewFilesystem()
	if err != nil {
		return err
	}
	if err := FileSystem.MountImage(me.local.Cluster); err != nil {
		return err
	}
	return nil
}
// step 7
func (me *LiteBuilder) ReMountImage() error {
	err := me.UnMountImage()
	if err != nil {
		return err
	}
	me.local.Cluster.Spec.Image = me.local.Config.ImageName
	return me.MountImage()
}
// step 8
func (me *LiteBuilder) InitDockerAndRegistry() error {
	mount := filepath.Join(common.DefaultClusterBaseDir(me.local.Cluster.Name), "mount")
	cmd := "cd %s  && chmod +x scripts/* && cd scripts && sh docker.sh && sh init-registry.sh 5000 %s"
	r, err := utils.CmdOutput("sh", "-c", fmt.Sprintf(cmd, mount, filepath.Join(mount, "registry")))
	if err != nil {
		logger.Error(fmt.Sprintf("Init docker and registry failed: %v", err))
		return err
	}
	logger.Info(string(r))
	return nil
}
// step 9
func (me *LiteBuilder) CacheImageToRegistry() error {
	var images []string
	var err error
	d := docker.Docker{}
	c := charts.Charts{}
	m := manifest.Manifests{}
	imageList := filepath.Join(common.DefaultClusterBaseDir(me.local.Cluster.Name), "mount", "manifests", "imageList")
	if utils.IsExist(imageList) {
		images, err = utils.ReadLines(imageList)
	}
	if helmImages, err := c.ListImages(me.local.Cluster.Name); err == nil {
		images = append(images, helmImages...)
	}
	if manifestImages, err := m.ListImages(me.local.Cluster.Name); err == nil {
		images = append(images, manifestImages...)
	}
	if err != nil {
		return err
	}
	d.ImagesPull(images)
	return nil
}
// step 10
func (me *LiteBuilder) AddUpperLayerToImage() error {
	var (
		err   error
		Image *v1.Image
	)
	m := filepath.Join(common.DefaultClusterBaseDir(me.local.Cluster.Name), "mount")
	err = mount.NewMountDriver().Unmount(m)
	if err != nil {
		return err
	}
	upper := filepath.Join(m, "upper")
	imageLayer := v1.Layer{
		Type:  "BASE",
		Value: "registry cache",
	}
	err = me.local.calculateLayerDigestAndPlaceIt(&imageLayer, upper)
	if err != nil {
		return err
	}
	Image, err = image.GetImageByName(me.local.Config.ImageName)
	if err != nil {
		return err
	}
	Image.Spec.Layers = append(Image.Spec.Layers, imageLayer)
	me.local.Image = Image
	err = me.local.updateImageIDAndSaveImage()
	if err != nil {
		return err
	}
	return nil
}
// step 11
func (me *LiteBuilder) Clear() error {
	return utils.CleanFiles(common.RawClusterfile, common.DefaultClusterBaseDir(me.local.Cluster.Name))
}

//#endregion

func (me *LiteBuilder) UnMountImage() error {
	var (
		FileSystem filesystem.Interface
		err        error
	)
	FileSystem, err = filesystem.NewFilesystem()
	if err != nil {
		logger.Warn(err)
		return err
	}
	return FileSystem.UnMountImage(me.local.Cluster)
}


func NewLiteBuilder(config *Config) (Interface, error) {
	localBuilder, err := NewLocalBuilder(config)
	if err != nil {
		return nil, err
	}
	return &LiteBuilder{
		local: localBuilder.(*LocalBuilder),
	}, nil
}
