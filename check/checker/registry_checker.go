// Copyright © 2021 Alibaba Group Holding Ltd.

package checker

import (
	"context"
	"fmt"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/image/distributionutil"
	"github.com/alibaba/sealer/utils"
)

type RegistryChecker struct {
	RegistryDomain string
}

func (r *RegistryChecker) Check() error {
	// check the existence of the docker.json ;
	authFile := common.DefaultRegistryAuthConfigDir()
	if !utils.IsFileExist(authFile) {
		return fmt.Errorf("registry auth info not found,please run 'sealer login' first")
	}
	// try to login with auth info
	authConfig, err := utils.GetDockerAuthInfoFromDocker(r.RegistryDomain)
	if err != nil {
		return fmt.Errorf("failed to get auth info, err: %s", err)
	}

	err = distributionutil.Login(context.Background(), &authConfig)
	if err != nil {
		return fmt.Errorf("%v authentication failed", r.RegistryDomain)
	}
	return nil
}

func NewRegistryChecker(registryDomain string) Checker {
	return &RegistryChecker{RegistryDomain: registryDomain}
}
