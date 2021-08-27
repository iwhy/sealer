// Copyright © 2021 Alibaba Group Holding Ltd.

package checker

// Checker: checking cluster status,such as node,pod,svc status.
type Checker interface {
	Check() error
}
