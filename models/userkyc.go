package models

const (
	//KYCNone used in UserKYCStruct.Tier to represent no KYC verification exists.
	KYCNone = 0
	//KYCBasic used in UserKYCStruct.Tier to represent basic verification only.
	KYCBasic = 1
	//KYCIntermediate used in UserKYCStruct.Tier to represent intermediate level of verification.
	KYCIntermediate = 2
	//KYCAdvanced used in UserKYCStruct.Tier to represent Advanced for Full level of verification.
	KYCAdvanced = 3
)

//UserKYCSignature Hold Verification status as well as who verified and when
type UserKYCSignature struct {
	Verified bool
	//UID of Person Node/Record who Verified
	Person uint64
	//The Time in which verification occured
	When uint64
}

//UserKYCTierOneStruct Basic Customer (User) information (Tier 1)
type UserKYCTierOneStruct struct {
	//Person is the UID of a Person Node/Record
	Person uint64
	//Address is the UID of a Address Node/Record
	Address uint64
	//Phone is the UID of a Phone Node/Record
	Phone uint64
	//NatID is a countries National Identification Number (SSN for USA) for the person
	NatID string
	//PictureOfIDPath is the Full Path and filename to the PNG containing a cropped scan of Picture Identification
	PictureOfIDPath string
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
}

//UserKYCTierTwoStruct Intermediate Customer (User) information (Tier 2)
type UserKYCTierTwoStruct struct {
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
	//TODO: Complete
}

//UserKYCTierThreeStruct Advanced Customer (User) information (Tier 3)
type UserKYCTierThreeStruct struct {
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
	//TODO: Complete
}

// UserKYCStruct contains users kyc "Know Your Customer" extended data
type UserKYCStruct struct {
	//User Current verified Tier Level (0..3)
	//
	Tier  uint
	Tier1 UserKYCTierOneStruct
	Tier2 UserKYCTierTwoStruct
	Tier3 UserKYCTierThreeStruct
}
