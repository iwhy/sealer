// Copyright © 2021 Alibaba Group Holding Ltd.

package runtime

const (
	AuditPolicyYml = "audit-policy.yml"
)

// static file should not be template, will never be changed while initialization
type StaticFile struct {
	DestinationDir string
	Name           string
}

//MasterStaticFiles Put static files here, can be moved to all master nodes before kubeadm execution
var MasterStaticFiles = []*StaticFile{
	{
		DestinationDir: "/etc/kubernetes",
		Name:           AuditPolicyYml,
	},
}
