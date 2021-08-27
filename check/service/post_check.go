// Copyright © 2021 Alibaba Group Holding Ltd.

package service

import (
	"github.com/alibaba/sealer/check/checker"
	"github.com/alibaba/sealer/logger"
)

type PostCheckerService struct {
}

func (d *PostCheckerService) Run() error {
	checkerList, err := d.init()
	if err != nil {
		logger.Error(err)
		return err
	}
	for _, checker := range checkerList {
		err = checker.Check()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *PostCheckerService) init() ([]checker.Checker, error) {
	var checkerList []checker.Checker
	checkerList = append(checkerList, checker.NewNodeChecker(), checker.NewPodChecker(), checker.NewSvcChecker())
	return checkerList, nil
}

func NewPostCheckerService() CheckerService {
	return &PostCheckerService{}
}
