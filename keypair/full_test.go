package keypair

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("keypair.Full", func() {
	var subject KeyPair

	JustBeforeEach(func() {
		subject = &Full{seed}
	})

	ItBehavesLikeAKP(&subject)

	type SignCase struct {
		Message   string
		Signature string
	}

	DescribeTable("Sign()",
		func(c SignCase) {
			sig, err := subject.Sign([]byte(c.Message))
			actual := hex.EncodeToString(sig)

			Expect(actual).To(Equal(c.Signature))
			Expect(err).To(BeNil())
		},

		Entry("hello", SignCase{
			"hello",
			"665c1a3d6a198a0c5ae63b6a1b4f6cd45f96f05f5f51d1040ef7f72fc390305d1a65cfb8057c87f484772aab757cb79720ee770cc28e58975157a6e7d2b3dc0e",
		}),
		Entry("this is a message", SignCase{
			"this is a message",
			"293d4a0bd309959c4d07b25777be1abba838e21b9e8cade56f4ff9ae8e2d55b385ae5dd2400854620b07c3cfff928f937aa1d789d8506e845b83122e6c66e50f",
		}),
	)

	Describe("LibP2PPrivKey()", func() {
		It("succeeds", func() {
			_, err := subject.LibP2PPrivKey()
			Expect(err).To(BeNil())
		})

	})

	Describe("LibP2PPubKey()", func() {
		It("succeeds", func() {
			_, err := subject.LibP2PPubKey()
			Expect(err).To(BeNil())
		})

	})
})
