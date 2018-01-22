package isolated

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("delete-buildpack command", func() {
	Context("when the stack is specified to disambiguate multiple buildpacks with the same name", func() {
		var dir string

		BeforeEach(func() {
			LoginCF()

			var err error
			dir, err = ioutil.TempDir("", "delete-buildpack-test")
			Expect(err).ToNot(HaveOccurred())

			filename := "manifest.yml"
			manifestFile := filepath.Join(dir, filename)
			err = ioutil.WriteFile(manifestFile, []byte("---\nstack:cflinuxfs2"), 0400)
			Expect(err).ToNot(HaveOccurred())

			session := CF("create-buildpack", "some-buildpack", dir, "1")
			Eventually(session).Should(Exit(0))
		})

		AfterEach(func() {
			Expect(os.RemoveAll(dir)).To(Succeed())
		})

		It("accepts stack argument and deletes the buildpack", func() {
			session := CF("delete-buildpack", "some-buildpack", "-s", "cflinuxfs2", "-f")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("OK"))
			Expect(session.Err).NotTo(Say("Incorrect Usage:"))
		})
	})
})
