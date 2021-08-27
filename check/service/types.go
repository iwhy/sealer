// Copyright © 2021 Alibaba Group Holding Ltd.

package service

type CheckerService interface {
	Run() error
}

type CheckerArgs struct {
	Pre  bool
	Post bool
}
