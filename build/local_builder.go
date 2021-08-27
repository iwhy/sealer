// Copyright © 2021 Alibaba Group Holding Ltd.

package build

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/alibaba/sealer/command"
	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/image/cache"
	"github.com/alibaba/sealer/image/reference"
	"github.com/alibaba/sealer/image/store"
	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/parser"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/archive"
	"github.com/alibaba/sealer/utils/mount"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

//由命令行获取的参数
type Config struct {
	BuildType string
	NoCache   bool
	ImageName string
}

//定义两个文件层，，overlayer2
type builderLayer struct {
	baseLayers []v1.Layer
	newLayers  []v1.Layer
}

// LocalBuilder: local builder using local provider to build a cluster image
// 主要包含Image相关的信息
// 注意这里buildLayer的写法
type LocalBuilder struct {
	Config               *Config
	Image                *v1.Image
	Cluster              *v1.Cluster
	ImageNamed           reference.Named
	ImageID              string
	Context              string
	KubeFileName         string
	LayerStore           store.LayerStore
	ImageStore           store.ImageStore
	ImageService         image.Service
	Prober               image.Prober
	FS                   store.Backend
	DockerImageStorePath string
	builderLayer
}

//#region build相关
func (me *LocalBuilder) Build(name string, context string, kubefileName string) error {
	err := me.initBuilder(name, context, kubefileName)
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

func (me *LocalBuilder) initBuilder(name string, context string, kubefileName string) error {
	named, err := reference.ParseToNamed(name)
	if err != nil {
		return err
	}

	me.ImageNamed = named
	me.Context = context
	me.KubeFileName = kubefileName
	return nil
}

func (me *LocalBuilder) GetBuildPipeLine() ([]func() error, error) {
	var buildPipeline []func() error
	if err := me.InitImageSpec(); err != nil {
		return nil, err
	}

	buildPipeline = append(buildPipeline,
		me.PullBaseImageNotExist,
		me.ExecBuild,
		me.UpdateImageMetadata)
	return buildPipeline, nil
}

// init default Image metadata
func (me *LocalBuilder) InitImageSpec() error {
	kubeFile, err := utils.ReadAll(me.KubeFileName)
	if err != nil {
		return fmt.Errorf("failed to load kubefile: %v", err)
	}
	me.Image = parser.NewParse().Parse(kubeFile)
	if me.Image == nil {
		return fmt.Errorf("failed to parse kubefile, image is nil")
	}

	layer0 := me.Image.Spec.Layers[0]
	if layer0.Type != common.FROMCOMMAND {
		return fmt.Errorf("first line of kubefile must start with FROM")
	}

	logger.Info("init image spec success!")
	return nil
}

func (me *LocalBuilder) PullBaseImageNotExist() (err error) {
	if me.Image.Spec.Layers[0].Value == common.ImageScratch {
		return nil
	}
	if err = me.ImageService.PullIfNotExist(me.Image.Spec.Layers[0].Value); err != nil {
		return fmt.Errorf("failed to pull baseImage: %v", err)
	}
	logger.Info("pull base image %s success", me.Image.Spec.Layers[0].Value)
	return nil
}

func (me *LocalBuilder) ExecBuild() error {
	err := me.updateBuilderLayers(me.Image)
	if err != nil {
		return err
	}

	var (
		canUseCache = !me.Config.NoCache
		parentID    = cache.ChainID("")
		newLayers   = me.newLayers
	)

	baseLayerPaths := getBaseLayersPath(me.baseLayers)
	chainSvc, err := cache.NewService()
	if err != nil {
		return err
	}

	for i := 0; i < len(newLayers); i++ {
		layer := &newLayers[i]
		logger.Info("run build layer: %s %s", layer.Type, layer.Value)
		if canUseCache {
			canUseCache, parentID = me.goCache(parentID, layer, chainSvc)
			// cache layer is empty layer
			if canUseCache {
				if layer.ID == "" {
					continue
				}
				baseLayerPaths = append(baseLayerPaths, me.FS.LayerDataDir(layer.ID))
				continue
			}
		}

		if layer.Type == common.COPYCOMMAND {
			err = me.execCopyLayer(layer)
			if err != nil {
				return err
			}
		} else {
			// exec other build cmd,need to mount
			err = me.execOtherLayer(layer, baseLayerPaths)
			if err != nil {
				return err
			}
		}

		if layer.ID == "" {
			continue
		}

		baseLayerPaths = append(baseLayerPaths, me.FS.LayerDataDir(layer.ID))
	}
	// todo need to collect docker images while build
	logger.Info("exec all build instructs success !")
	return nil
}

//#region ExecBuild
func (me *LocalBuilder) updateBuilderLayers(image *v1.Image) error {
	// we do not check the len of layers here, because we checked it before.
	// remove the first layer of image
	var (
		layer0    = image.Spec.Layers[0]
		baseImage *v1.Image
		err       error
	)

	// and the layer 0 must be from layer
	if layer0.Value == common.ImageScratch {
		// give a empty image
		baseImage = &v1.Image{}
	} else {
		baseImage, err = me.ImageStore.GetByName(image.Spec.Layers[0].Value)
		if err != nil {
			return fmt.Errorf("failed to get base image while updating base layers, err: %s", err)
		}
	}

	me.baseLayers = append([]v1.Layer{}, baseImage.Spec.Layers...)
	me.newLayers = append([]v1.Layer{}, image.Spec.Layers[1:]...)
	if len(me.baseLayers)+len(me.newLayers) > maxLayerDeep {
		return errors.New("current number of layers exceeds 128 layers")
	}
	return nil
}
// used in build stage, where the image still has from layer
func getBaseLayersPath(layers []v1.Layer) (res []string) {
	for _, layer := range layers {
		if layer.ID != "" {
			res = append(res, filepath.Join(common.DefaultLayerDir, layer.ID.Hex()))
		}
	}
	return res
}


// run COPY command, because user can overwrite some file like Cluster file, or build a base image
func (me *LocalBuilder) execCopyLayer(layer *v1.Layer) error {
	//count layer hash;create layer dir ;update image layer hash
	tempDir, err := utils.MkTmpdir()
	if err != nil {
		return fmt.Errorf("failed to create %s:%v", tempDir, err)
	}
	defer utils.CleanDir(tempDir)

	err = me.execLayer(layer, tempDir)
	if err != nil {
		return fmt.Errorf("failed to exec layer %v:%v", layer, err)
	}

	if err = me.calculateLayerDigestAndPlaceIt(layer, tempDir); err != nil {
		return err
	}

	if err = me.SetCacheID(layer); err != nil {
		return err
	}

	return nil
}

func (me *LocalBuilder) execOtherLayer(layer *v1.Layer, lowLayers []string) error {
	tempTarget, err := utils.MkTmpdir()
	if err != nil {
		return fmt.Errorf("failed to create %s:%v", tempTarget, err)
	}
	tempUpper, err := utils.MkTmpdir()
	if err != nil {
		return fmt.Errorf("failed to create %s:%v", tempUpper, err)
	}
	defer utils.CleanDirs(tempTarget, tempUpper)

	if err = me.mountAndExecLayer(layer, tempTarget, tempUpper, lowLayers...); err != nil {
		return err
	}
	if err = me.calculateLayerDigestAndPlaceIt(layer, tempUpper); err != nil {
		return err
	}
	return nil
}


func (me *LocalBuilder) execLayer(layer *v1.Layer, tempTarget string) error {
	// exec layer cmd;
	if layer.Type == common.COPYCOMMAND {
		src := filepath.Join(me.Context, strings.Fields(layer.Value)[0])
		dest := ""
		if utils.IsDir(src) {
			// src is dir
			dest = filepath.Join(tempTarget, strings.Fields(layer.Value)[1], filepath.Base(src))
		} else {
			// src is file
			dest = filepath.Join(tempTarget, strings.Fields(layer.Value)[1], strings.Fields(layer.Value)[0])
		}
		return utils.RecursionCopy(src, dest)
	}
	if layer.Type == common.RUNCOMMAND || layer.Type == common.CMDCOMMAND {
		cmd := fmt.Sprintf(common.CdAndExecCmd, tempTarget, layer.Value)
		output, err := command.NewSimpleCommand(cmd).Exec()
		logger.Info(output)
		if err != nil {
			if me.Config.BuildType == common.LiteBuild {
				logger.Warn(fmt.Sprintf("failed to exec %s, err: %v", cmd, err))
				return nil
			}
			return fmt.Errorf("failed to exec %s, err: %v", cmd, err)
		}
	}
	return nil
}


//#endregion ExecBuild

func (me *LocalBuilder) UpdateImageMetadata() error {
	me.setClusterFileToImage()
	me.squashBaseImageLayerIntoCurrentImage()
	err := me.updateImageIDAndSaveImage()
	if err != nil {
		return fmt.Errorf("failed to updateImageIDAndSaveImage, err: %v", err)
	}

	logger.Info("update image %s to image metadata success !", me.ImageNamed.Raw())
	return nil
}

//#region UpdateImageMetadata()

//#region step 1 setClusterFileToImage

// setClusterFileToImage: set cluster file whatever build type is
func (me *LocalBuilder) setClusterFileToImage() {
	clusterFileData := GetRawClusterFile(me.Image)
	me.addImageAnnotations(common.ImageAnnotationForClusterfile, clusterFileData)
}
// GetClusterFile from user build context or from base image
func GetRawClusterFile(im *v1.Image) string {
if im.Spec.Layers[0].Value == common.ImageScratch {
data, err := ioutil.ReadFile(filepath.Join("etc", common.DefaultClusterFileName))
if err != nil {
return ""
}
return string(data)
}
// find cluster file from context
if clusterFile := getClusterFileFromContext(im); clusterFile != nil {
logger.Info("get cluster file from context success!")
return string(clusterFile)
}
// find cluster file from base image
clusterFile := image.GetClusterFileFromImage(im.Spec.Layers[0].Value)
if clusterFile != "" {
logger.Info("get cluster file from base image success!")
return clusterFile
}
return ""
}
func getClusterFileFromContext(image *v1.Image) []byte {
	for i := range image.Spec.Layers {
		layer := image.Spec.Layers[i]
		if layer.Type == common.COPYCOMMAND && strings.Fields(layer.Value)[0] == common.DefaultClusterFileName {
			if clusterFile, _ := utils.ReadAll(strings.Fields(layer.Value)[0]); clusterFile != nil {
				return clusterFile
			}
		}
	}
	return nil
}

// GetClusterFile from user build context or from base image
func (me *LocalBuilder) addImageAnnotations(key, value string) {
	if me.Image.Annotations == nil {
		me.Image.Annotations = make(map[string]string)
	}
	me.Image.Annotations[key] = value
}
//#endregion

//#region step 2 squashBaseImageLayerIntoCurrentImage

func (me *LocalBuilder) squashBaseImageLayerIntoCurrentImage() {
	me.Image.Spec.Layers = append(me.baseLayers, me.newLayers...)
}

//#endregion

//#region step 3 updateImageIDAndSaveImage

func (me *LocalBuilder) updateImageIDAndSaveImage() error {
	imageID, err := generateImageID(*me.Image)
	if err != nil {
		return err
	}

	me.Image.Spec.ID = imageID
	return me.ImageStore.Save(*me.Image, me.ImageNamed.Raw())
}

func generateImageID(image v1.Image) (string, error) {
	imageBytes, err := yaml.Marshal(image)
	if err != nil {
		return "", err
	}
	imageID := digest.FromBytes(imageBytes).Hex()
	return imageID, nil
}


//#endregion

//#endregion

//#endregion

//#region 公开了5个函数，，但是其它地方并未引用到

//This function only has meaning for copy layers
func (me *LocalBuilder) SetCacheID(layer *v1.Layer) error {
	baseDir := me.Context
	layerDgst, _, err := archive.TarCanonicalDigest(filepath.Join(baseDir, strings.Fields(layer.Value)[0]))
	if err != nil {
		return err
	}

	return me.FS.SetMetadata(layer.ID, cacheID, []byte(layerDgst.String()))
}

func (me *LocalBuilder) mountAndExecLayer(layer *v1.Layer, tempTarget, tempUpper string, lowLayers ...string) error {
	driver := mount.NewMountDriver()
	err := driver.Mount(tempTarget, tempUpper, lowLayers...)
	if err != nil {
		return fmt.Errorf("failed to mount target %s:%v", tempTarget, err)
	}
	defer func() {
		if err = driver.Unmount(tempTarget); err != nil {
			logger.Warn(fmt.Errorf("failed to umount %s:%v", tempTarget, err))
		}
	}()

	err = me.execLayer(layer, tempTarget)
	if err != nil {
		return fmt.Errorf("failed to exec layer %v:%v", layer, err)
	}
	return nil
}

func (me *LocalBuilder) calculateLayerDigestAndPlaceIt(layer *v1.Layer, tempTarget string) error {
	layerDigest, err := me.LayerStore.RegisterLayerForBuilder(tempTarget)
	if err != nil {
		return fmt.Errorf("failed to register layer, err: %v", err)
	}

	layer.ID = layerDigest
	return nil
}

func (me *LocalBuilder) goCache(parentID cache.ChainID, layer *v1.Layer, cacheService cache.Service) (continueCache bool, chainID cache.ChainID) {
	var (
		srcDigest = digest.Digest("")
		err       error
	)

	// specially for copy command, we would generate digest of src file as srcDigest.
	// and use srcDigest as cacheID to generate a cacheLayer, eventually use the cacheLayer
	// to hit the cache layer
	if layer.Type == common.COPYCOMMAND {
		srcDigest, err = me.generateSourceFilesDigest(strings.Fields(layer.Value)[0])
		if err != nil {
			logger.Warn("failed to generate src digest, discard cache, err: %s", err)
		}
	}

	cacheLayer := cacheService.NewCacheLayer(*layer, srcDigest)
	cacheLayerID, err := me.Prober.Probe(parentID.String(), &cacheLayer)
	if err != nil {
		logger.Debug("failed to probe cache for %+v, err: %s", layer, err)
		return false, ""
	}
	// cache hit
	logger.Info("---> Using cache %v", cacheLayerID)
	layer.ID = cacheLayerID
	cID, err := cacheLayer.ChainID(parentID)
	if err != nil {
		return false, ""
	}
	return true, cID
}

func (me *LocalBuilder) generateSourceFilesDigest(path string) (digest.Digest, error) {
	baseDir := me.Context
	layerDgst, _, err := archive.TarCanonicalDigest(filepath.Join(baseDir, path))
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return layerDgst, nil
}

//#endregion

//实现工厂的new方法，返回接口类型
// 就是初始化一个LocalBuilder struct
func NewLocalBuilder(config *Config) (Interface, error) {
	layerStore, err := store.NewDefaultLayerStore()
	if err != nil {
		return nil, err
	}

	imageStore, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}

	service, err := image.NewImageService()
	if err != nil {
		return nil, err
	}

	fs, err := store.NewFSStoreBackend()
	if err != nil {
		return nil, fmt.Errorf("failed to init store backend, err: %s", err)
	}

	prober := image.NewImageProber(service, config.NoCache)

	dockerImageStorePath, err := utils.MkTmpdir()
	if err != nil {
		return nil, fmt.Errorf("failed to create %s:%v", dockerImageStorePath, err)
	}

	return &LocalBuilder{
		Config:               config,
		LayerStore:           layerStore,
		ImageStore:           imageStore,
		ImageService:         service,
		Prober:               prober,
		FS:                   fs,
		DockerImageStorePath: dockerImageStorePath,
		builderLayer: builderLayer{
			// for skip golang ci
			baseLayers: []v1.Layer{},
			newLayers:  []v1.Layer{},
		},
	}, nil
}
