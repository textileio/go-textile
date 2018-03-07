package models

type IdentityType string
const (
	phoneNumber IdentityType = "phone_number"
	emailAddress IdentityType = "email_address"
)

type Identity struct {
	Type IdentityType `bson:"type" json:"type" binding:"required"`
	Value string `bson:"value" json:"value" binding:"required"`
}
