package keypair

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("keypair.FromAddress", func() {
	var subject KeyPair

	JustBeforeEach(func() {
		subject = &FromAddress{address}
	})

	ItBehavesLikeAKP(&subject)

	Describe("Sign()", func() {
		It("fails", func() {
			_, err := subject.Sign(message)
			Expect(err).To(HaveOccurred())
		})

	})

	Describe("LibP2PPrivKey()", func() {
		It("fails", func() {
			_, err := subject.LibP2PPrivKey()
			Expect(err).To(HaveOccurred())
		})

	})

	Describe("LibP2PPubKey()", func() {
		It("succeeds", func() {
			_, err := subject.LibP2PPubKey()
			Expect(err).To(BeNil())
		})

	})
})
