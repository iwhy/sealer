// Copyright © 2021 Alibaba Group Holding Ltd.

package main

import (
	"github.com/alibaba/sealer/sealer/boot"
	"github.com/alibaba/sealer/sealer/cmd"
)

func main() {
	err := boot.OnBoot()
	if err != nil {
		panic(err)
	}
	cmd.Execute()
}
