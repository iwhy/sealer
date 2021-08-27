// Copyright © 2021 Alibaba Group Holding Ltd.

package test

import (
	"github.com/alibaba/sealer/test/suites/image"

	"github.com/alibaba/sealer/test/suites/registry"
	"github.com/alibaba/sealer/test/testhelper/settings"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("sealer login", func() {
	Context("login docker registry", func() {
		AfterEach(func() {
			registry.Logout()
		})
		It("with correct name and password", func() {
			image.CheckLoginResult(
				settings.RegistryURL,
				settings.RegistryUsername,
				settings.RegistryPasswd,
				true)
		})
		It("with incorrect name and password", func() {
			image.CheckLoginResult(
				settings.RegistryURL,
				settings.RegistryPasswd,
				settings.RegistryUsername,
				false)
		})
		It("with only name", func() {
			image.CheckLoginResult(
				settings.RegistryURL,
				settings.RegistryUsername,
				"",
				false)
		})
		It("with only password", func() {
			image.CheckLoginResult(
				settings.RegistryURL,
				"",
				settings.RegistryPasswd,
				false)
		})
		It("with only registryURL", func() {
			image.CheckLoginResult(
				settings.RegistryURL,
				"",
				"",
				false)
		})
	})
})
