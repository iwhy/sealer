// Copyright © 2021 Alibaba Group Holding Ltd.

package runtime

import (
	"fmt"
	"path/filepath"

	"github.com/alibaba/sealer/logger"

	"github.com/alibaba/sealer/utils"
)

type RegistryConfig struct {
	IP     string `yaml:"ip,omitempty"`
	Domain string `yaml:"domain,omitempty"`
	Port   string `yaml:"port,omitempty"`
}

func getRegistryHost(rootfs, defaultRegistry string) (host string) {
	cf := GetRegistryConfig(rootfs, defaultRegistry)
	return fmt.Sprintf("%s %s", cf.IP, cf.Domain)
}

func GetRegistryConfig(rootfs, defaultRegistry string) *RegistryConfig {
	var config RegistryConfig
	var DefaultConfig = &RegistryConfig{
		IP:     utils.GetHostIP(defaultRegistry),
		Domain: SeaHub,
		Port:   "5000",
	}
	registryConfigPath := filepath.Join(rootfs, "/etc/registry.yaml")
	if !utils.IsFileExist(registryConfigPath) {
		logger.Warn(fmt.Sprintf("failed to find %s", registryConfigPath))
		return DefaultConfig
	}
	err := utils.UnmarshalYamlFile(registryConfigPath, &config)
	logger.Info(fmt.Sprintf("show registry info, IP: %s, Domain: %s", config.IP, config.Domain))
	if err != nil {
		logger.Error("Failed to read registry config! ")
		return DefaultConfig
	}
	if config.IP == "" {
		config.IP = DefaultConfig.IP
	} else {
		config.IP = utils.GetHostIP(config.IP)
	}
	if config.Port == "" {
		config.Port = DefaultConfig.Port
	}
	if config.Domain == "" {
		config.Domain = DefaultConfig.Domain
	}
	return &config
}

const registryName = "sealer-registry"

//Only use this for join and init, due to the initiation operations
func (d *Default) EnsureRegistry() error {
	cf := GetRegistryConfig(d.Rootfs, d.Masters[0])
	cmd := fmt.Sprintf("cd %s/scripts && sh init-registry.sh %s %s/registry", d.Rootfs, cf.Port, d.Rootfs)
	return d.SSH.CmdAsync(cf.IP, cmd)
}

func (d *Default) RecycleRegistry() error {
	cf := GetRegistryConfig(d.Rootfs, d.Masters[0])
	cmd := fmt.Sprintf("docker stop %s || true && docker rm %s || true", registryName, registryName)
	return d.SSH.CmdAsync(cf.IP, cmd)
}
