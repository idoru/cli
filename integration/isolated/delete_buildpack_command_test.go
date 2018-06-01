package isolated

import (
	. "code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("delete-buildpack command", func() {
	BeforeEach(func() {
		LoginCF()
	})

	Context("there is exactly one buildpack with the specified name", func() {
		BeforeEach(func() {
			session := CF("create-buildpack", "some-buildpack-1", "../assets/test_buildpacks/simple_buildpack-cflinuxfs2-v1.0.0.zip", "1")
			Eventually(session).Should(Exit(0))
		})

		It("deletes the specified buildpack", func() {
			By("passing the associated stack")
			session := CF("delete-buildpack", "some-buildpack-1", "-s", "cflinuxfs2", "-f")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("OK"))

			session = CF("create-buildpack", "some-buildpack-1", "../assets/test_buildpacks/simple_buildpack-cflinuxfs2-v1.0.0.zip", "1")
			Eventually(session).Should(Exit(0))

			By("passing no stack")
			session = CF("delete-buildpack", "some-buildpack-1", "-f")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("OK"))
		})
	})

	Context("there is a second buildpack with the same name", func() {
		BeforeEach(func() {
			session := CF("create-buildpack", "some-buildpack-2", "../assets/test_buildpacks/simple_buildpack-cflinuxfs2-v1.0.0.zip", "1")
			Eventually(session).Should(Exit(0))

			session = CF("create-buildpack", "some-buildpack-2", "../assets/test_buildpacks/simple_buildpack-windows2012R2-v1.0.0.zip", "1")
			Eventually(session).Should(Exit(0))
		})

		It("properly handles ambiguity", func() {
			By("failing when no stack specified")
			session := CF("delete-buildpack", "some-buildpack-2", "-f")
			Eventually(session).Should(Exit(1))
			Eventually(session.Out).Should(Say("FAILED"))

			By("deleting the buildpack when the associated stack is specified")
			session = CF("delete-buildpack", "some-buildpack-2", "-s", "cflinuxfs2", "-f")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("OK"))

			session = CF("delete-buildpack", "some-buildpack-2", "-s", "windows2012R2", "-f")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("OK"))
		})
	})
})
