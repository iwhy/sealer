// Copyright © 2021 Alibaba Group Holding Ltd.

package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/alibaba/sealer/common"
	dockerstreams "github.com/docker/cli/cli/streams"
	dockerjsonmessage "github.com/docker/docker/pkg/jsonmessage"

	"github.com/alibaba/sealer/logger"
	"github.com/alibaba/sealer/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Docker struct {
	Username string
	Password string
}

func (d Docker) ImagesPull(images []string) {
	for _, image := range utils.RemoveDuplicate(images) {
		if image == "" {
			continue
		}
		if err := d.ImagePull(image); err != nil {
			logger.Warn(fmt.Sprintf("Image %s pull failed: %v", image, err))
		}
	}
}

// ImagePull函数重载，，核心还是调用"github.com/docker/docker/client"的ImagePull
func (d Docker) ImagePull(image string) error {
	var ImagePullOptions types.ImagePullOptions
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	//获取docker授权
	if d.Username != "" && d.Password != "" {
		authConfig := types.AuthConfig{
			Username: d.Username,
			Password: d.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		ImagePullOptions = types.ImagePullOptions{RegistryAuth: authStr}
	}

	out, err := cli.ImagePull(ctx, image, ImagePullOptions)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	err = dockerjsonmessage.DisplayJSONMessagesToStream(out, dockerstreams.NewOut(common.StdOut), nil)
	if err != nil && err != io.ErrClosedPipe {
		logger.Warn("error occurs in display progressing, err: %s", err)
	}
	return nil
}


func (d Docker) ImagesPullByImageListFile(fileName string) {
	data, err := utils.ReadLines(fileName)
	if err != nil {
		logger.Error(fmt.Sprintf("Read image list failed: %v", err))
	}
	d.ImagesPull(data)
}

func (d Docker) ImagesPullByList(images []string) {
	d.ImagesPull(images)
}

