// Copyright © 2021 Alibaba Group Holding Ltd.

package checker

import "fmt"

type PodNotReadyError struct {
	name string
}

func (e *PodNotReadyError) Error() string {
	return fmt.Sprintf("pod  %s is not ready", e.name)
}

type NotFindReadyTypeError struct {
	name string
}

func (e *NotFindReadyTypeError) Error() string {
	return fmt.Sprintf("pod %s has't Ready Type", e.name)
}
