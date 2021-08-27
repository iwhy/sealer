// Copyright © 2021 Alibaba Group Holding Ltd.

package store

import (
	"github.com/alibaba/sealer/common"
	pkgutils "github.com/alibaba/sealer/utils"

	"path/filepath"
)

func DeleteImageLocal(imageID string) (err error) {
	return deleteImage(imageID)
}

func deleteImage(imageID string) error {
	file := filepath.Join(common.DefaultImageDBRootDir, imageID+common.YamlSuffix)
	if pkgutils.IsFileExist(file) {
		err := pkgutils.CleanFiles(file)
		if err != nil {
			return err
		}
	}
	return nil
}
