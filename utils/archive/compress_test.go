// Copyright © 2021 Alibaba Group Holding Ltd.

//nolint
package archive

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const basePath = "/tmp"

const fileContent = "content"

type fileDef struct {
	name    string
	content string
}

type dirDef struct {
	path   string
	files  []fileDef
	subDir []dirDef
}

var filesToCreate = []dirDef{
	{
		path: "testDirA",
		files: []fileDef{
			{
				name:    "testFileA",
				content: fileContent,
			},
			{
				name:    "testFileB",
				content: fileContent,
			},
		},
		subDir: []dirDef{
			{
				path: "testDirC",
				files: []fileDef{
					{
						name:    "testFileA",
						content: fileContent,
					},
					{
						name:    "testFileB",
						content: fileContent,
					},
				},
			},
		},
	},
	{
		path: "testDirB",
		files: []fileDef{
			{
				name:    "testFileA",
				content: fileContent,
			},
			{
				name:    "testFileB",
				content: fileContent,
			},
		},
	},
}

func makeDir(root string, d dirDef) error {
	currentDir := filepath.Join(root, d.path)
	err := os.MkdirAll(currentDir, 0755)
	if err != nil {
		return err
	}

	for _, file := range d.files {
		_, err = os.Create(filepath.Join(currentDir, file.name))
		if err != nil {
			return err
		}
	}

	for _, sub := range d.subDir {
		err = makeDir(currentDir, sub)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestTarWithoutRootDir(t *testing.T) {
	digest, _, err := TarCanonicalDigest("/Users/eric/Workspace/src/sealer/empty")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(digest)
}

func TestTarWithRootDir(t *testing.T) {
	reader, err := TarWithRootDir("./hash.go")
	if err != nil {
		t.Error(err)
	}

	tmp, err := ioutil.TempFile("/tmp", "tar")
	_, err = io.Copy(tmp, reader)
	if err != nil {
		t.Error(err)
	}
}
