// Copyright © 2021 Alibaba Group Holding Ltd.

package checker

import (
	"text/template"

	"github.com/alibaba/sealer/common"

	corev1 "k8s.io/api/core/v1"

	"github.com/alibaba/sealer/client"
	"github.com/alibaba/sealer/logger"
)

const (
	ReadyNodeStatus    = "Ready"
	NotReadyNodeStatus = "NotReady"
)

type NodeChecker struct {
}

type NodeClusterStatus struct {
	ReadyCount       uint32
	NotReadyCount    uint32
	NodeCount        uint32
	NotReadyNodeList []string
}

func (n *NodeChecker) Check() error {
	// check if all the node is ready
	c, err := client.NewClientSet()
	if err != nil {
		logger.Info("failed to create k8s client  %v", err)
		return nil
	}
	nodes, err := client.ListNodes(c)
	if err != nil {
		return err
	}
	var notReadyNodeList []string
	var readyCount uint32 = 0
	var nodeCount uint32
	var notReadyCount uint32 = 0
	for _, node := range nodes.Items {
		nodeIP, nodePhase := GetNodeStatus(&node)
		if nodePhase != ReadyNodeStatus {
			notReadyCount++
			notReadyNodeList = append(notReadyNodeList, nodeIP)
		} else {
			readyCount++
		}
	}
	nodeCount = notReadyCount + readyCount
	nodeClusterStatus := NodeClusterStatus{
		ReadyCount:       readyCount,
		NotReadyCount:    notReadyCount,
		NodeCount:        nodeCount,
		NotReadyNodeList: notReadyNodeList,
	}
	err = n.Output(nodeClusterStatus)
	if err != nil {
		return err
	}
	return nil
}

func (n *NodeChecker) Output(nodeCLusterStatus NodeClusterStatus) error {
	//t1, err := template.ParseFiles("templates/node_checker.tpl")
	t := template.New("node_checker")
	t, err := t.Parse(
		`Cluster Node Status
  ReadyNode: {{ .ReadyCount }}/{{ .NodeCount }}
  {{ if (gt .NotReadyCount 0 ) -}}
  Not Ready Node List:
    {{- range .NotReadyNodeList }}
    NodeIP: {{ . }}
    {{- end }}
  {{ end }}
`)
	if err != nil {
		panic(err)
	}
	t = template.Must(t, err)
	err = t.Execute(common.StdOut, nodeCLusterStatus)
	if err != nil {
		logger.Error("node checker template can not excute %s", err)
		return err
	}
	return nil
}

func GetNodeStatus(node *corev1.Node) (IP string, Phase string) {
	if len(node.Status.Addresses) < 1 {
		return "", ""
	}
	for _, address := range node.Status.Addresses {
		if address.Type == "InternalIP" {
			IP = address.Address
		}
	}
	if IP == "" {
		IP = node.Status.Addresses[0].Address
	}
	Phase = NotReadyNodeStatus
	for _, condition := range node.Status.Conditions {
		if condition.Type == ReadyNodeStatus {
			if condition.Status == "True" {
				Phase = ReadyNodeStatus
			} else {
				Phase = NotReadyNodeStatus
			}
		}
	}
	return IP, Phase
}

func NewNodeChecker() Checker {
	return &NodeChecker{}
}
