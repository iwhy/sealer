// Copyright © 2021 Alibaba Group Holding Ltd.

package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alibaba/sealer/image/store"

	"github.com/alibaba/sealer/runtime"

	infraUtils "github.com/alibaba/sealer/infra/utils"

	"github.com/alibaba/sealer/utils"

	"github.com/pkg/errors"

	"github.com/alibaba/sealer/logger"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/image"

	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils/mount"
	"github.com/alibaba/sealer/utils/ssh"
)

const (
	RemoteChmod = "cd %s  && chmod +x scripts/* && cd scripts && sh init.sh"
)

type Interface interface {
	MountRootfs(cluster *v1.Cluster, hosts []string) error
	UnMountRootfs(cluster *v1.Cluster) error
	MountImage(cluster *v1.Cluster) error
	UnMountImage(cluster *v1.Cluster) error
	Clean(cluster *v1.Cluster) error
}

type FileSystem struct {
	imageStore store.ImageStore
}

func (c *FileSystem) Clean(cluster *v1.Cluster) error {
	return utils.CleanFiles(common.GetClusterWorkDir(cluster.Name), common.DefaultClusterBaseDir(cluster.Name), common.DefaultKubeConfigDir())
}

func (c *FileSystem) umountImage(cluster *v1.Cluster) error {
	mountdir := common.DefaultMountCloudImageDir(cluster.Name)
	if utils.IsFileExist(mountdir) {
		var err error
		err = infraUtils.Retry(10, time.Second, func() error {
			err = mount.NewMountDriver().Unmount(mountdir)
			if err != nil {
				return err
			}
			return os.RemoveAll(mountdir)
		})
		if err != nil {
			logger.Warn("failed to unmount dir %s,err: %v", mountdir, err)
		}
	}
	return nil
}

func (c *FileSystem) mountImage(cluster *v1.Cluster) error {
	mountdir := common.DefaultMountCloudImageDir(cluster.Name)
	upperDir := filepath.Join(mountdir, "upper")
	if utils.IsDir(mountdir) {
		if utils.IsFileExist(upperDir) {
			utils.CleanDir(upperDir)
		} else {
			logger.Info("image already mounted")
			return nil
		}
	}
	//get layers
	Image, err := c.imageStore.GetByName(cluster.Spec.Image)
	if err != nil {
		return err
	}
	layers, err := image.GetImageLayerDirs(Image)
	if err != nil {
		return fmt.Errorf("get layers failed: %v", err)
	}
	driver := mount.NewMountDriver()
	if err = os.MkdirAll(upperDir, 0744); err != nil {
		return fmt.Errorf("create upperdir failed, %s", err)
	}
	if err = driver.Mount(mountdir, upperDir, layers...); err != nil {
		return fmt.Errorf("mount files failed %v", err)
	}
	return nil
}

func (c *FileSystem) MountImage(cluster *v1.Cluster) error {
	err := c.mountImage(cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *FileSystem) UnMountImage(cluster *v1.Cluster) error {
	err := c.umountImage(cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *FileSystem) MountRootfs(cluster *v1.Cluster, hosts []string) error {
	clusterRootfsDir := common.DefaultTheClusterRootfsDir(cluster.Name)
	//scp roofs to all Masters and Nodes,then do init.sh
	if err := mountRootfs(hosts, clusterRootfsDir, cluster); err != nil {
		return fmt.Errorf("mount rootfs failed %v", err)
	}
	return nil
}

func (c *FileSystem) UnMountRootfs(cluster *v1.Cluster) error {
	//do clean.sh,then remove all Masters and Nodes roofs
	IPList := append(cluster.Spec.Masters.IPList, cluster.Spec.Nodes.IPList...)
	config := runtime.GetRegistryConfig(common.DefaultTheClusterRootfsDir(cluster.Name), cluster.Spec.Masters.IPList[0])
	if utils.NotIn(config.IP, IPList) {
		IPList = append(IPList, config.IP)
	}
	if err := unmountRootfs(IPList, cluster); err != nil {
		return err
	}
	return nil
}

func mountRootfs(ipList []string, target string, cluster *v1.Cluster) error {
	SSH := ssh.NewSSHByCluster(cluster)
	config := runtime.GetRegistryConfig(
		common.DefaultTheClusterRootfsDir(cluster.Name),
		cluster.Spec.Masters.IPList[0])
	if err := ssh.WaitSSHReady(SSH, ipList...); err != nil {
		return errors.Wrap(err, "check for node ssh service time out")
	}
	var wg sync.WaitGroup
	var flag bool
	var mutex sync.Mutex
	src := common.DefaultMountCloudImageDir(cluster.Name)
	localHostAddrs, err := utils.IsLocalHostAddrs()
	if err != nil {
		return err
	}
	// TODO scp sdk has change file mod bug
	initCmd := fmt.Sprintf(RemoteChmod, target)
	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			err = CopyFiles(SSH, ip == config.IP, utils.IsLocalIP(ip, localHostAddrs), ip, src, target)
			if err != nil {
				logger.Error("copy rootfs failed %v", err)
				mutex.Lock()
				flag = true
				mutex.Unlock()
			}
			err = SSH.CmdAsync(ip, initCmd)
			if err != nil {
				logger.Error("exec init.sh failed %v", err)
				mutex.Lock()
				flag = true
				mutex.Unlock()
			}
		}(ip)
	}
	wg.Wait()
	if flag {
		return fmt.Errorf("mountRootfs failed")
	}
	return nil
}

func CopyFiles(ssh ssh.Interface, isRegistry bool, isLocal bool, ip, src, target string) error {
	logger.Info(fmt.Sprintf(" %s the local host: %t ,", ip, isLocal))
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to copy files %s", err)
	}
	if isLocal {
		if isRegistry {
			return utils.RecursionCopy(src, target)
		}
		for _, f := range files {
			if f.Name() == common.RegistryDirName {
				continue
			}
			err = utils.RecursionCopy(filepath.Join(src, f.Name()), filepath.Join(target, f.Name()))
			if err != nil {
				return fmt.Errorf("failed to local copy sub files %v", err)
			}
		}
	} else {
		if isRegistry {
			return ssh.Copy(ip, src, target)
		}
		for _, f := range files {
			if f.Name() == common.RegistryDirName {
				continue
			}
			err = ssh.Copy(ip, filepath.Join(src, f.Name()), filepath.Join(target, f.Name()))
			if err != nil {
				return fmt.Errorf("failed to copy sub files %v", err)
			}
		}
	}
	return nil
}

func unmountRootfs(ipList []string, cluster *v1.Cluster) error {
	SSH := ssh.NewSSHByCluster(cluster)
	var wg sync.WaitGroup
	var flag bool
	var mutex sync.Mutex
	clusterRootfsDir := common.DefaultTheClusterRootfsDir(cluster.Name)
	execClean := fmt.Sprintf("/bin/sh -c "+common.DefaultClusterClearFile, cluster.Name)
	rmRootfs := fmt.Sprintf("rm -rf %s", clusterRootfsDir)
	for _, ip := range ipList {
		wg.Add(1)
		go func(IP string) {
			defer wg.Done()
			if err := SSH.CmdAsync(IP, execClean, rmRootfs); err != nil {
				logger.Error("%s:exec %s failed, %s", IP, execClean, err)
				mutex.Lock()
				flag = true
				mutex.Unlock()
				return
			}
		}(ip)
	}
	wg.Wait()
	if flag {
		return fmt.Errorf("unmountRootfs failed")
	}
	return nil
}

func NewFilesystem() (Interface, error) {
	dis, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}

	return &FileSystem{imageStore: dis}, nil
}
