// Copyright © 2021 Alibaba Group Holding Ltd.

package store

import (
	"io/ioutil"
	"path/filepath"
)

//var supportedDigestAlgo = map[string]bool{
//	digest.SHA256.String(): true,
//	digest.SHA384.String(): true,
//	digest.SHA512.String(): true,
//}

func getDirListInDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, file := range files {
		// avoid adding some other dirs created by users
		if file.IsDir() {
			dirs = append(dirs, filepath.Join(dir, file.Name()))
		}
	}
	return dirs, nil
}

func traverseLayerDB(layerDBRoot string) ([]string, error) {
	// TODO maybe there no need to traverse layerdb, just clarify how many sha supported in a list
	shaDirs, err := getDirListInDir(layerDBRoot)
	if err != nil {
		return nil, err
	}

	var layerDirs []string
	for _, shaDir := range shaDirs {
		layerDirList, err := getDirListInDir(shaDir)
		if err != nil {
			return nil, err
		}
		layerDirs = append(layerDirs, layerDirList...)
	}
	return layerDirs, nil
}
