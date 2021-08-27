// Copyright © 2021 Alibaba Group Holding Ltd.

package service

import "fmt"

type PreCheckerService struct {
}

func (d *PreCheckerService) Run() error {
	fmt.Println("Pre check cluster")
	return nil
}

func NewPreCheckerService() CheckerService {
	return &PreCheckerService{}
}
