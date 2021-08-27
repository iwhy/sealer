// Copyright © 2021 Alibaba Group Holding Ltd.

package checker

import (
	"text/template"

	"github.com/alibaba/sealer/common"

	corev1 "k8s.io/api/core/v1"

	"github.com/alibaba/sealer/client"
	"github.com/alibaba/sealer/logger"
)

type SvcChecker struct {
}
type SvcNamespaceStatus struct {
	NamespaceName       string
	ServiceCount        int
	EndpointCount       int
	UnhealthServiceList []string
}

type SvcClusterStatus struct {
	SvcNamespaceStatusList []*SvcNamespaceStatus
}

func (n *SvcChecker) Check() error {
	// check if all the svc is ok
	c, err := client.NewClientSet()
	if err != nil {
		logger.Info("failed to create k8s client  %v", err)
		return nil
	}
	namespaceSvcList, err := client.ListAllNamespacesSvcs(c)
	var svcNamespaceStatusList []*SvcNamespaceStatus
	if err != nil {
		return err
	}
	for _, svcNamespace := range namespaceSvcList {
		serviceCount := len(svcNamespace.ServiceList.Items)
		var unhaelthService []string
		var endpointCount = 0
		endpointsList, err := client.GetEndpointsList(c, svcNamespace.Namespace.Name)
		if err != nil {
			break
		}
		for _, service := range svcNamespace.ServiceList.Items {
			if IsExistEndpoint(endpointsList, service.Name) {
				endpointCount++
			} else {
				unhaelthService = append(unhaelthService, service.Name)
			}
		}
		svcNamespaceStatus := SvcNamespaceStatus{
			NamespaceName:       svcNamespace.Namespace.Name,
			ServiceCount:        serviceCount,
			EndpointCount:       endpointCount,
			UnhealthServiceList: unhaelthService,
		}
		svcNamespaceStatusList = append(svcNamespaceStatusList, &svcNamespaceStatus)
	}
	err = n.Output(svcNamespaceStatusList)
	if err != nil {
		return err
	}
	return nil
}

func (n *SvcChecker) Output(svcNamespaceStatusList []*SvcNamespaceStatus) error {
	t := template.New("svc_checker")
	t, err := t.Parse(
		`Cluster Service Status
  {{- range . }}
  Namespace: {{ .NamespaceName }}
  HealthService: {{ .EndpointCount }}/{{ .ServiceCount }}
  UnhealthServiceList:
    {{- range .UnhealthServiceList }}
    ServiceName: {{ . }}
    {{- end }}
  {{- end }}
`)
	if err != nil {
		panic(err)
	}
	t = template.Must(t, err)
	err = t.Execute(common.StdOut, svcNamespaceStatusList)
	if err != nil {
		logger.Error("service checker template can not excute %s", err)
		return err
	}
	return nil
}

func IsExistEndpoint(endpointList *corev1.EndpointsList, serviceName string) bool {
	for _, ep := range endpointList.Items {
		if ep.Name == serviceName {
			if len(ep.Subsets) > 0 {
				return true
			}
		}
	}
	return false
}

func NewSvcChecker() Checker {
	return &SvcChecker{}
}
