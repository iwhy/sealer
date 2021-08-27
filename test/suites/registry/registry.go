// Copyright © 2021 Alibaba Group Holding Ltd.

package registry

import (
	"fmt"

	"github.com/alibaba/sealer/test/testhelper"
	"github.com/alibaba/sealer/test/testhelper/settings"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func Login() {
	sess, err := testhelper.Start(fmt.Sprintf("%s login %s -u %s -p %s", settings.DefaultSealerBin, settings.RegistryURL,
		settings.RegistryUsername,
		settings.RegistryPasswd))

	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(sess).Should(gbytes.Say(fmt.Sprintf("login %s success", settings.RegistryURL)))
	gomega.Eventually(sess, settings.MaxWaiteTime).Should(gexec.Exit(0))
}

func Logout() {
	testhelper.DeleteFileLocally(DefaultRegistryAuthConfigDir())
}

// DefaultRegistryAuthConfigDir using root privilege to run sealer cmd at e2e test
func DefaultRegistryAuthConfigDir() string {
	return settings.DefaultRegistryAuthFileDir
}
