// Copyright © 2021 Alibaba Group Holding Ltd.

package runtime

import (
	"fmt"
	"sync"

	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
)

func (d *Default) reset(cluster *v1.Cluster) error {
	err := d.resetNodes(cluster.Spec.Nodes.IPList)
	if err != nil {
		logger.Error("failed to clean nodes %v", err)
	}
	err = d.resetMasters(cluster.Spec.Masters.IPList)
	if err != nil {
		logger.Error("failed to clean masters %v", err)
	}
	return d.RecycleRegistry()
}
func (d *Default) resetNodes(nodes []string) error {
	if len(nodes) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			if err := d.resetNode(node); err != nil {
				logger.Error("delete node %s failed %v", node, err)
			}
		}(node)
	}
	wg.Wait()

	return nil
}
func (d *Default) resetMasters(nodes []string) error {
	if len(nodes) == 0 {
		return nil
	}
	for _, node := range nodes {
		if err := d.resetNode(node); err != nil {
			logger.Error("delete master %s failed %v", node, err)
		}
	}
	return nil
}
func (d *Default) resetNode(node string) error {
	host := utils.GetHostIP(node)
	if err := d.SSH.CmdAsync(host, fmt.Sprintf(RemoteCleanMasterOrNode, vlogToStr(d.Vlog)),
		fmt.Sprintf(RemoteRemoveAPIServerEtcHost, d.APIServer),
		fmt.Sprintf(RemoteRemoveAPIServerEtcHost, getRegistryHost(d.Rootfs, d.Masters[0]))); err != nil {
		return err
	}
	return nil
}
