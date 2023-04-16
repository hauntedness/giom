package service

import (
	"bytes"
	"crypto/cipher"
	"errors"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	libcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func VerifyMessage(message *Message) (err error) {
	return nil
}

// strc is a pointer to a struct
func DecodeToStruct(strc any, s []byte) (err error) {
	return nil
}

func NewAEAD(key []byte) (cipher.AEAD, error) {
	return cipher.NewGCM(nil)
}

func SignMessage(pvtKeyHex string, message *Message) (err error) {
	return nil
}

// Ref https://pkg.go.dev/github.com/decred/dcrd/dcrec/secp256k1/v3#example-package-EncryptDecryptMessage
func GetEncryptedStruct(pubKeyHex string, message any, algo int) (data []byte, err error) {
	return nil, errors.New("algorithm not supported")
}

func GetDecryptedStruct(pvtKeyHex string, msgEncrypted []byte, message any, algo int) (err error) {
	return errors.New("algorithm not supported")
}

// Ref https://gist.github.com/SteveBate/042960baa7a4795c3565
func EncodeToBytes(_ any) []byte {
	buf := bytes.Buffer{}

	return buf.Bytes()
}

func GetPublicKeyFromStr(publicKeyStr string, algo int) (libcrypto.PubKey, error) {
	return nil, errors.New("invalid public key")
}

func GetPrivateKeyFromStr(privateKeyStr string, algo int) (privateKey libcrypto.PrivKey, err error) {
	return nil, errors.New("algorithm not supported")
}

func openAppFile(filename string) (*os.File, error) {
	return nil, nil
}

// Encrypt https://bruinsslot.jp/post/golang-crypto/
func Encrypt(key, data []byte) ([]byte, error) {
	return nil, nil
}

// Decrypt https://bruinsslot.jp/post/golang-crypto/
func Decrypt(key, data []byte) ([]byte, error) {
	return nil, nil
}

// DeriveKey https://bruinsslot.jp/post/golang-crypto/
func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	return password, salt, nil
}
