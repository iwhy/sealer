// Copyright © 2021 Alibaba Group Holding Ltd.

package apply

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/filesystem"
	"github.com/alibaba/sealer/guest"
	"github.com/alibaba/sealer/image"
	"github.com/alibaba/sealer/infra"
	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/runtime"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/ssh"
)

const ApplyCluster = "chmod +x %s && %s apply -f %s"

type CloudApplier struct {
	*DefaultApplier
}

func NewAliCloudProvider(cluster *v1.Cluster) (Interface, error) {
	imgService, err := image.NewImageService()
	if err != nil {
		return nil, err
	}

	fs, err := filesystem.NewFilesystem()
	if err != nil {
		return nil, err
	}

	gs, err := guest.NewGuestManager()
	if err != nil {
		return nil, err
	}

	d := &DefaultApplier{
		ClusterDesired: cluster,
		ImageManager:   imgService,
		FileSystem:     fs,
		Guest:          gs,
	}
	return &CloudApplier{d}, nil
}

func (c *CloudApplier) ScaleDownNodes(cluster *v1.Cluster) (isScaleDown bool, err error) {
	if cluster == nil {
		return false, nil
	}
	logger.Info("desired master %s, current master %s, desired nodes %s, current nodes %s", c.ClusterDesired.Spec.Masters.Count,
		cluster.Spec.Masters.Count,
		c.ClusterDesired.Spec.Nodes.Count,
		cluster.Spec.Nodes.Count)
	if c.ClusterDesired.Spec.Masters.Count >= cluster.Spec.Masters.Count &&
		c.ClusterDesired.Spec.Nodes.Count >= cluster.Spec.Nodes.Count {
		return false, nil
	}

	MastersToJoin, MastersToDelete := utils.GetDiffHosts(cluster.Spec.Masters, c.ClusterDesired.Spec.Masters)
	NodesToJoin, NodesToDelete := utils.GetDiffHosts(cluster.Spec.Nodes, c.ClusterDesired.Spec.Nodes)
	if len(MastersToJoin) != 0 || len(NodesToJoin) != 0 {
		return false, fmt.Errorf("should not scale up and down at same time")
	}

	err = DeleteNodes(append(MastersToDelete, NodesToDelete...))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CloudApplier) Apply() error {
	var err error
	cluster := c.ClusterDesired
	clusterCurrent, err := GetCurrentCluster()
	if err != nil {
		return fmt.Errorf("failed to get current cluster %v", err)
	}

	cloudProvider, err := infra.NewDefaultProvider(cluster)
	if err != nil {
		return err
	}
	if cloudProvider == nil {
		return fmt.Errorf("new cloud provider failed")
	}
	err = cloudProvider.Apply()
	if err != nil {
		return fmt.Errorf("apply infra failed %v", err)
	}
	if cluster.DeletionTimestamp != nil {
		return nil
	}
	err = utils.SaveClusterfile(cluster)
	if err != nil {
		return err
	}

	scaleDown, err := c.ScaleDownNodes(clusterCurrent)
	if err != nil {
		return fmt.Errorf("failed to scale down nodes %v", err)
	}
	if scaleDown {
		//  infra already delete the host, if continue to apply will not found the host and
		//  return ssh error
		logger.Info("scale the cluster success")
		return nil
	}

	client, err := ssh.NewSSHClientWithCluster(cluster)
	if err != nil {
		return fmt.Errorf("prepare cluster ssh client failed %v", err)
	}

	cluster.Spec.Provider = common.BAREMETAL
	err = utils.MarshalYamlToFile(common.TmpClusterfile, cluster)
	if err != nil {
		return fmt.Errorf("marshal tmp cluster file failed %v", err)
	}
	defer func() {
		if err := utils.CleanFiles(common.TmpClusterfile); err != nil {
			logger.Error("failed to clean %s, err: %v", common.TmpClusterfile, err)
		}
	}()

	err = runtime.PreInitMaster0(client.SSH, client.Host)
	if err != nil {
		return err
	}

	err = client.SSH.CmdAsync(client.Host, fmt.Sprintf(ApplyCluster, common.RemoteSealerPath, common.RemoteSealerPath, common.TmpClusterfile))
	if err != nil {
		return err
	}

	err = runtime.GetKubectlAndKubeconfig(client.SSH, client.Host)
	if err != nil {
		return fmt.Errorf("failed to copy kubeconfig and kubectl %v", err)
	}

	return nil
}

func (c *CloudApplier) Delete() error {
	t := metav1.Now()
	c.ClusterDesired.DeletionTimestamp = &t
	host := c.ClusterDesired.GetAnnotationsByKey(common.Eip)
	err := c.Apply()
	if err != nil {
		return err
	}
	if err := utils.RemoveFileContent(common.EtcHosts, fmt.Sprintf("%s %s", host, common.APIServerDomain)); err != nil {
		logger.Warn(err)
	}

	if err := utils.CleanFiles(common.DefaultKubeConfigDir(), common.GetClusterWorkDir(c.ClusterDesired.Name), common.TmpClusterfile, common.KubectlPath); err != nil {
		logger.Warn(err)
		return nil
	}

	return nil
}
