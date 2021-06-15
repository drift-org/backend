package helpers_test

import (
	"github.com/drift-org/backend/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Challenge helper", func ()  {
	Describe("Find Challenge", func() {
		It("Basic", helpers.TestWithMongo("FindChallenge-Basic", func() {
			Expect(true).To(BeTrue())
		}))
	})
})