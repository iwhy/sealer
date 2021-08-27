// Copyright © 2021 Alibaba Group Holding Ltd.

package service

import "fmt"

type DefaultCheckerService struct {
}

func (d *DefaultCheckerService) Run() error {
	fmt.Println("check cluster")
	return nil
}

func NewDefaultCheckerService() CheckerService {
	return &DefaultCheckerService{}
}
