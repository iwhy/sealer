// Copyright © 2021 Alibaba Group Holding Ltd.

package archive

import (
	"fmt"
	"io"

	"github.com/opencontainers/go-digest"
)

const emptySHA256TarDigest = "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"

func TarCanonicalDigest(path string) (digest.Digest, int64, error) {
	tarReader, err := TarWithoutRootDir(path)
	if err != nil {
		return "", 0, fmt.Errorf("unable to tar on %s, err: %s", path, err)
	}
	defer tarReader.Close()

	digester := digest.Canonical.Digester()
	size, err := io.Copy(digester.Hash(), tarReader)
	if err != nil {
		return "", 0, err
	}
	layerDigest := digester.Digest()
	if layerDigest == emptySHA256TarDigest {
		return "", 0, nil
	}

	return layerDigest, size, nil
}
