package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("loki-proxy", func() {
	const timeout, poll = 3 * time.Minute, 5 * time.Second

	It("routes an annotated namespace to its own tenant only", func() {
		Eventually(lokiNamespaces("solar"), timeout, poll).Should(ConsistOf("solar"))
		Eventually(lokiNamespaces("wind"), timeout, poll).Should(ConsistOf("wind"))
	})

	It("routes a namespace without annotation to the default tenant", func() {
		Eventually(lokiNamespaces("platform"), timeout, poll).Should(ContainElement("green"))
		Expect(lokiNamespaces("platform")()).NotTo(ContainElements("solar", "wind"))
	})
})
