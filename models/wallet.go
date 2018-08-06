package models

import (
	"crypto/ecdsa"
)

// Wallet stores private and public keys
type Wallet struct {
	UserID     uint64 //AccountStruct|dGraph Node UID? Who Owns this wallet...For authentication
	label      string //Name for wallet (Limit length to 32 )
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}
