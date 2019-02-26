package keypair

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestBuild(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Package: github.com/textileio/go-textile/keypair")
}

var (
	address   = "P4WysQUVHg2GDMctx7Sf3DjLRAfRZqsZoCPr6aec7JgWpSRq"
	seed      = "STqdutoFiPV9Nzh96L7bvSZNzrHkQiB2phDW1ySG9Na1kkGE"
	id        = "12D3KooWBRG69PMx6EakNbkp4bQRqGqsHou5eau1MgnkzZDmzezU"
	hint      = [4]byte{0x55, 0x79, 0x7f, 0xe1}
	message   = []byte("hello")
	signature = []byte{
		0x66, 0x5c, 0x1a, 0x3d, 0x6a, 0x19, 0x8a, 0x0c, 0x5a, 0xe6, 0x3b, 0x6a,
		0x1b, 0x4f, 0x6c, 0xd4, 0x5f, 0x96, 0xf0, 0x5f, 0x5f, 0x51, 0xd1, 0x04,
		0x0e, 0xf7, 0xf7, 0x2f, 0xc3, 0x90, 0x30, 0x5d, 0x1a, 0x65, 0xcf, 0xb8,
		0x05, 0x7c, 0x87, 0xf4, 0x84, 0x77, 0x2a, 0xab, 0x75, 0x7c, 0xb7, 0x97,
		0x20, 0xee, 0x77, 0x0c, 0xc2, 0x8e, 0x58, 0x97, 0x51, 0x57, 0xa6, 0xe7,
		0xd2, 0xb3, 0xdc, 0x0e,
	}
)

func ItBehavesLikeAKP(subject *KeyPair) {

	// NOTE: subject will only be valid to dereference when inside am "It"
	// example.

	Describe("Address()", func() {
		It("returns the correct address", func() {
			Expect((*subject).Address()).To(Equal(address))
		})
	})

	Describe("Id()", func() {
		It("returns the correct id", func() {
			pid, err := (*subject).Id()
			if err != nil {
				panic(err)
			}
			Expect(pid.Pretty()).To(Equal(id))
		})
	})

	Describe("Hint()", func() {
		It("returns the correct hint", func() {
			Expect((*subject).Hint()).To(Equal(hint))
		})
	})

	type VerifyCase struct {
		Message   []byte
		Signature []byte
		Case      types.GomegaMatcher
	}

	DescribeTable("Verify()",
		func(vc VerifyCase) {
			Expect((*subject).Verify(vc.Message, vc.Signature)).To(vc.Case)
		},
		Entry("correct", VerifyCase{message, signature, BeNil()}),
		Entry("empty signature", VerifyCase{message, []byte{}, Equal(ErrInvalidSignature)}),
		Entry("empty message", VerifyCase{[]byte{}, signature, Equal(ErrInvalidSignature)}),
		Entry("different message", VerifyCase{[]byte("diff"), signature, Equal(ErrInvalidSignature)}),
		Entry("malformed signature", VerifyCase{message, signature[0:10], Equal(ErrInvalidSignature)}),
	)
}

type ParseCase struct {
	Input    string
	TypeCase types.GomegaMatcher
	ErrCase  types.GomegaMatcher
}

var _ = DescribeTable("keypair.Parse()",
	func(c ParseCase) {
		kp, err := Parse(c.Input)

		Expect(kp).To(c.TypeCase)
		Expect(err).To(c.ErrCase)
	},

	Entry("a valid address", ParseCase{
		Input:    "P6NKHguuQnGraNM58WZTh2tPmRAnQQn1YzHQZYPmRt8WABDF",
		TypeCase: BeAssignableToTypeOf(&FromAddress{}),
		ErrCase:  BeNil(),
	}),
	Entry("a corrupted address", ParseCase{
		Input:    "P6NKHguuQnGraNM58WZTh2tPmRAnQQn1YzHQZYPmRt8WALT5",
		TypeCase: BeNil(),
		ErrCase:  HaveOccurred(),
	}),
	Entry("a valid seed", ParseCase{
		Input:    "SV8k5RKcUg1ZtN6qvcUGXfvqjv2nGeBhxy7sUnG1AxM5jB23",
		TypeCase: BeAssignableToTypeOf(&Full{}),
		ErrCase:  BeNil(),
	}),
	Entry("a corrupted seed", ParseCase{
		Input:    "ST1JhsR34aiso4N5RgiM2F7XRSxihACsrTDSm64VfnGTfeop",
		TypeCase: BeNil(),
		ErrCase:  HaveOccurred(),
	}),
	Entry("a blank string", ParseCase{
		Input:    "",
		TypeCase: BeNil(),
		ErrCase:  HaveOccurred(),
	}),
)
