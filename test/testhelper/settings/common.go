// Copyright © 2021 Alibaba Group Holding Ltd.

package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

const (
	SealerBinPath                     = "/usr/local/bin/sealer"
	ImageName                         = "sealer_test_image_"
	DefaultImageDomain                = "registry.cn-qingdao.aliyuncs.com"
	DefaultImageRepo                  = "sealer-io"
	DefaultImageName                  = "kubernetes:v1.19.9"
	DefaultRegistryAuthFileDir        = "/root/.docker"
	DefaultClusterFileNeedToBeCleaned = "/root/.sealer/%s/Clusterfile"
	SubCmdBuildOfSealer               = "build"
	SubCmdApplyOfSealer               = "apply"
	SubCmdDeleteOfSealer              = "delete"
	SubCmdRunOfSealer                 = "run"
	SubCmdLoginOfSealer               = "login"
	SubCmdTagOfSealer                 = "tag"
	SubCmdPullOfSealer                = "pull"
	SubCmdListOfSealer                = "images"
	SubCmdPushOfSealer                = "push"
	SubCmdRmiOfSealer                 = "rmi"
	SubCmdForceRmiOfSealer            = "frmi"
	DefaultSSHPassword                = "Sealer123"
	ImageAnnotationForClusterfile     = "sea.aliyun.com/ClusterFile"
)

const (
	FileMode0755 = 0755
	FileMode0644 = 0644
)
const (
	LocalBuild = "local"
)
const (
	BAREMETAL         = "BAREMETAL"
	AliCloud          = "ALI_CLOUD"
	CONTAINER         = "CONTAINER"
	DefaultImage      = "registry.cn-qingdao.aliyuncs.com/sealer-io/kubernetes:v1.19.9"
	ClusterNameForRun = "my-cluster"
	TMPClusterFile    = "/tmp/Clusterfile"
	ClusterWorkDir    = "/root/.sealer/%s"
)

var (
	DefaultPollingInterval time.Duration
	MaxWaiteTime           time.Duration
	DefaultWaiteTime       time.Duration
	DefaultSealerBin       = ""
	DefaultTestEnvDir      = ""
	RegistryURL            = os.Getenv("REGISTRY_URL")
	RegistryUsername       = os.Getenv("REGISTRY_USERNAME")
	RegistryPasswd         = os.Getenv("REGISTRY_PASSWORD")
	CustomImageName        = os.Getenv("IMAGE_NAME")

	AccessKey     = os.Getenv("ACCESSKEYID")
	AccessSecret  = os.Getenv("ACCESSKEYSECRET")
	Region        = os.Getenv("RegionID")
	TestImageName = "" //default: registry.cn-qingdao.aliyuncs.com/sealer-io/kubernetes:v1.19.9
)

func GetClusterWorkDir(clusterName string) string {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Sprintf(ClusterWorkDir, clusterName)
	}
	return filepath.Join(home, ".sealer", clusterName)
}

func GetClusterWorkClusterfile(clusterName string) string {
	return filepath.Join(GetClusterWorkDir(clusterName), "Clusterfile")
}
