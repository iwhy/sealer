// Copyright © 2021 Alibaba Group Holding Ltd.

package utils

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/alibaba/sealer/common"
	v1 "github.com/alibaba/sealer/types/api/v1"

	"sigs.k8s.io/yaml"
)

func UnmarshalYamlFile(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file %s to %s", file, reflect.TypeOf(obj))
	}
	return nil
}

func MarshalYamlToFile(file string, obj interface{}) error {
	//Marshal把对象转为内存中的data
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	//WriteFile保存未指定路径的文件
	if err = WriteFile(file, data); err != nil {
		return err
	}
	return nil
}

func SaveClusterfile(cluster *v1.Cluster) error {
	fileName := common.GetClusterWorkClusterfile(cluster.Name)
	err := MkFileFullPathDir(fileName)
	if err != nil {
		return fmt.Errorf("mkdir failed %s %v", fileName, err)
	}
	err = MarshalYamlToFile(fileName, cluster)
	if err != nil {
		return fmt.Errorf("marshal cluster file failed %v", err)
	}
	return nil
}
