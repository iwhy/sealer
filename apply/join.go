// Copyright © 2021 Alibaba Group Holding Ltd.

package apply

import (
	"strconv"
	"strings"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
	"github.com/alibaba/sealer/utils"
)

func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		logger.Error("String to digit conversion failed:", err)
		return 0
	}
	return num
}

func removeIPListDuplicatesAndEmpty(ipList []string) []string {
	count := len(ipList)
	var newList []string
	for i := 0; i < count; i++ {
		if (i > 0 && ipList[i-1] == ipList[i]) || len(ipList[i]) == 0 {
			continue
		}
		newList = append(newList, ipList[i])
	}
	return newList
}

func JoinApplierFromArgs(clusterfile string, joinArgs *common.RunArgs) Interface {
	cluster := &v1.Cluster{}
	if err := utils.UnmarshalYamlFile(clusterfile, cluster); err != nil {
		logger.Error("clusterfile parsing failed, please check:", err)
		return nil
	}
	if joinArgs.Nodes == "" && joinArgs.Masters == "" {
		logger.Error("The node or master parameter was not committed")
		return nil
	}
	if cluster.Spec.Provider == "BAREMETAL" {
		if IsIPList(joinArgs.Nodes) || IsIPList(joinArgs.Masters) {
			margeMasters := append(cluster.Spec.Masters.IPList, strings.Split(joinArgs.Masters, ",")...)
			margeNodes := append(cluster.Spec.Nodes.IPList, strings.Split(joinArgs.Nodes, ",")...)
			cluster.Spec.Masters.IPList = removeIPListDuplicatesAndEmpty(margeMasters)
			cluster.Spec.Nodes.IPList = removeIPListDuplicatesAndEmpty(margeNodes)
		} else {
			logger.Error("Parameter error:", "The current mode should submit iplist！")
			return nil
		}
	} else if IsNumber(joinArgs.Nodes) || IsNumber(joinArgs.Masters) {
		cluster.Spec.Masters.Count = strconv.Itoa(StrToInt(cluster.Spec.Masters.Count) + StrToInt(joinArgs.Masters))
		cluster.Spec.Nodes.Count = strconv.Itoa(StrToInt(cluster.Spec.Nodes.Count) + StrToInt(joinArgs.Nodes))
	} else {
		logger.Error("Parameter error:", "The number of join masters or nodes that must be submitted to use cloud service！")
		return nil
	}
	if err := utils.MarshalYamlToFile(clusterfile, cluster); err != nil {
		logger.Error("clusterfile save failed, please check:", err)
		return nil
	}

	applier, err := NewApplier(cluster)
	if err != nil {
		logger.Error("failed to init applier, err: %s", err)
		return nil
	}
	return applier
}
