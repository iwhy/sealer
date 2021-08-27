// Copyright © 2021 Alibaba Group Holding Ltd.

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func (d Docker) DockerRmi(imageID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	if _, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{Force: true, PruneChildren: true}); err != nil {
		return err
	}

	//client.

	return nil
}

func (d Docker) ImagesList() ([]*types.ImageSummary, error) {
	var List []*types.ImageSummary
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		List = append(List, &image)
	}
	return List, nil
}
