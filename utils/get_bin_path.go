// Copyright © 2021 Alibaba Group Holding Ltd.

package utils

import (
	"os"
	"path/filepath"

	"github.com/alibaba/sealer/common"
)

func ExecutableFilePath() string {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, common.ExecBinaryFileName)
}
