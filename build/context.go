// Copyright © 2021 Alibaba Group Holding Ltd.

package build

import (
	"fmt"
	"os"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/runtime"
	"github.com/alibaba/sealer/utils"
)

//sendBuildContext:send local build context to remote server
func (c *CloudBuilder) sendBuildContext() (err error) {
	// if remote cluster already exist,no need to pre init master0
	if !c.SSH.IsFileExist(c.RemoteHostIP, common.RemoteSealerPath) {
		err = runtime.PreInitMaster0(c.SSH, c.RemoteHostIP)
		if err != nil {
			return fmt.Errorf("failed to prepare cluster env %v", err)
		}
	}

	// tar local build context
	tarFileName := fmt.Sprintf(common.TmpTarFile, utils.GenUniqueID(32))
	if _, isExist := utils.CheckCmdIsExist("tar"); !isExist {
		return fmt.Errorf("local server muster support tar cmd")
	}
	context := fmt.Sprintf("%s %s", c.local.KubeFileName, c.local.Context)
	if _, err := utils.RunSimpleCmd(fmt.Sprintf(common.ZipCmd, tarFileName, context)); err != nil {
		return fmt.Errorf("failed to create context file: %v", err)
	}
	defer func() {
		if err = os.Remove(tarFileName); err != nil {
			logger.Warn("failed to cleanup local temp file %s:%v", tarFileName, err)
		}
	}()
	// send to remote server
	workdir := fmt.Sprintf(common.DefaultWorkDir, c.local.Cluster.Name)
	if err = c.SSH.Copy(c.RemoteHostIP, tarFileName, tarFileName); err != nil {
		return fmt.Errorf("failed to copy tar file: %s, err: %v", tarFileName, err)
	}
	// unzip remote context
	err = c.SSH.CmdAsync(c.RemoteHostIP, fmt.Sprintf(common.UnzipCmd, workdir, tarFileName, workdir))
	if err != nil {
		return err
	}
	logger.Info("send build context to %s success !", c.RemoteHostIP)
	return nil
}

func (c *CloudBuilder) changeBuilderContext() {
	c.local.Context = "."
}
