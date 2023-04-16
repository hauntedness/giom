package service

import (
	"time"
)

type Account struct {
	UpdatedAt time.Time
	// PrivateKeyEnc is encrypted private key(private key itself is represented in hex string)
	PrivateKeyEnc string
	// PublicKey is hex representation of public key as string
	PublicKey    string `gorm:"primaryKey"`
	PublicImage  []byte
	PrivateImage []byte
}

func (a *Account) PrivateKey(passwd string) (string, error) {
	return "", nil
}
