package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cortex-proxy", func() {
	const timeout, poll = 3 * time.Minute, 5 * time.Second

	It("routes an annotated namespace to its own tenant only", func() {
		Eventually(mimirNamespaces("solar"), timeout, poll).Should(ConsistOf("solar"))
		Eventually(mimirNamespaces("wind"), timeout, poll).Should(ConsistOf("wind"))
	})

	It("routes a namespace without annotation to the default tenant", func() {
		Eventually(mimirNamespaces("platform"), timeout, poll).Should(ContainElement("green"))
		Expect(mimirNamespaces("platform")()).NotTo(ContainElements("solar", "wind"))
	})
})
