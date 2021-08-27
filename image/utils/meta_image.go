// Copyright © 2021 Alibaba Group Holding Ltd.

package utils

import (
	"strings"

	"github.com/alibaba/sealer/image/store"

	"github.com/spf13/cobra"
)

func SimilarImageList(imageArg string) (similarImageList []string, err error) {
	is, err := store.NewDefaultImageStore()
	if err != nil {
		return nil, err
	}

	metadataMap, err := is.GetImageMetadataMap()
	if err != nil {
		return nil, err
	}
	for _, imageMetadata := range metadataMap {
		imageMeta := imageMetadata
		if !strings.Contains(imageMeta.Name, imageArg) && imageArg != "" {
			continue
		}
		similarImageList = append(similarImageList, imageMeta.Name)
	}
	return similarImageList, nil
}

func ImageListFuncForCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	similarImages, err := SimilarImageList(toComplete)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	return similarImages, cobra.ShellCompDirectiveNoFileComp
}
