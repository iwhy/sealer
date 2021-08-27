// Copyright © 2021 Alibaba Group Holding Ltd.

package client

import (
	"context"
	"path/filepath"

	"github.com/alibaba/sealer/common"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type NamespacePod struct {
	Namespace v1.Namespace
	PodList   *v1.PodList
}

type NamespaceSvc struct {
	Namespace   v1.Namespace
	ServiceList *v1.ServiceList
}

func NewClientSet() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(common.DefaultKubeConfigDir(), "config")
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "new kube build config failed")
	}

	return kubernetes.NewForConfig(config)
}

func ListNodes(client *kubernetes.Clientset) (*v1.NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "get cluster nodes failed")
	}
	return nodes, nil
}

func DeleteNode(client *kubernetes.Clientset, name string) error {
	err := client.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "delete cluster nodes failed")
	}
	return nil
}

func ListNamespaces(client *kubernetes.Clientset) (*v1.NamespaceList, error) {
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get namespace")
	}
	return namespaceList, nil
}

func ListAllNamespacesPods(client *kubernetes.Clientset) ([]*NamespacePod, error) {
	namespaceList, err := ListNamespaces(client)
	if err != nil {
		return nil, err
	}
	var namespacePodList []*NamespacePod
	for _, ns := range namespaceList.Items {
		pods, err := client.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "get all namespace pods failed ")
		}
		namespacePod := NamespacePod{
			Namespace: ns,
			PodList:   pods,
		}
		namespacePodList = append(namespacePodList, &namespacePod)
	}

	return namespacePodList, nil
}

func ListAllNamespacesSvcs(client *kubernetes.Clientset) ([]*NamespaceSvc, error) {
	namespaceList, err := ListNamespaces(client)
	if err != nil {
		return nil, err
	}
	var namespaceSvcList []*NamespaceSvc
	for _, ns := range namespaceList.Items {
		svcs, err := client.CoreV1().Services(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "get all namespace pods failed ")
		}
		namespaceSvc := NamespaceSvc{
			Namespace:   ns,
			ServiceList: svcs,
		}
		namespaceSvcList = append(namespaceSvcList, &namespaceSvc)
	}
	return namespaceSvcList, nil
}

func GetEndpointsList(client *kubernetes.Clientset, namespace string) (*v1.EndpointsList, error) {
	endpointsList, err := client.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "get the endpoint in the %s namespace", namespace)
	}
	return endpointsList, nil
}

func ListSvcs(client *kubernetes.Clientset, namespace string) (*v1.ServiceList, error) {
	svcs, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "get all namespace pods failed ")
	}
	return svcs, nil
}
