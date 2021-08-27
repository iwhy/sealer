// Copyright © 2021 Alibaba Group Holding Ltd.

package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alibaba/sealer/utils"

	"github.com/opencontainers/go-digest"
)

type DistributionMetadataItem struct {
	SourceRepository   string        `json:"source_repository"`
	LayerContentDigest digest.Digest `json:"layer_content_digest"`
}

// distributionMetadata is the data from {layerdb}/distribution_layer_digest
// which indicate that digest of compressedlayerStream in specific registry and repository
type DistributionMetadata []DistributionMetadataItem

func (fs *filesystem) LoadDistributionMetadata(layerID LayerID) (map[string]digest.Digest, error) {
	var (
		layerDBPath = fs.LayerDBDir(layerID.ToDigest())
		metadatas   = DistributionMetadata{}
		res         = map[string]digest.Digest{}
	)
	distributionMetadataFile, err := os.Open(filepath.Join(layerDBPath, "distribution_layer_digest"))
	if err != nil {
		return res, nil
	}
	defer distributionMetadataFile.Close()

	err = json.NewDecoder(distributionMetadataFile).Decode(&metadatas)
	if err != nil {
		return res, err
	}

	for _, item := range metadatas {
		res[item.SourceRepository] = item.LayerContentDigest
	}

	return res, nil
}

func (fs *filesystem) addDistributionMetadata(layerID LayerID, newMetadatas map[string]digest.Digest) error {
	// load from distribution_layer_digest
	metadataMap, err := fs.LoadDistributionMetadata(layerID)
	if err != nil {
		return err
	}
	// override metadata items, and add new metadata
	for key, value := range newMetadatas {
		metadataMap[key] = value
	}

	distributionMetadatas := DistributionMetadata{}
	for key, value := range metadataMap {
		distributionMetadatas = append(distributionMetadatas, DistributionMetadataItem{
			SourceRepository:   key,
			LayerContentDigest: value,
		})
	}

	distributionMetadatasJSON, err := json.Marshal(&distributionMetadatas)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(fs.LayerDBDir(layerID.ToDigest()), "distribution_layer_digest"), distributionMetadatasJSON)
}
