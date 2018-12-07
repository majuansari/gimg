// aes256.go
// Author Andrey Izman <izmanw@gmail.com>
// Copyright 2018 Andrey Izman
// License MIT

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	b64 "encoding/base64"
)

// Encrypts text with the passphrase
func Encrypt(text, passphrase, saltText string) string {
	salt := saltText

	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	pad := __PKCS5Padding([]byte(text), block.BlockSize())
	ecb := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(pad))
	ecb.CryptBlocks(encrypted, pad)

	return b64.RawURLEncoding.EncodeToString([]byte(string(salt) + string(encrypted)))
}

// Decrypts encrypted text with the passphrase
func Decrypt(encrypted, passphrase, saltText string) string {
	ct, _ := b64.RawURLEncoding.DecodeString(encrypted)
	if len(ct) < 16 || string(ct[:16]) != saltText {
		return ""
	}

	salt := ct[0:16]
	ct = ct[16:]
	key, iv := __DeriveKeyAndIv(passphrase, string(salt))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	dst := make([]byte, len(ct))
	cbc.CryptBlocks(dst, ct)

	return string(__PKCS5Trimming(dst))
}

func __PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func __PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func __DeriveKeyAndIv(passphrase string, salt string) (string, string) {
	salted := ""
	dI := ""

	for len(salted) < 48 {
		md := md5.New()
		md.Write([]byte(dI + passphrase + salt))
		dM := md.Sum(nil)
		dI = string(dM[:16])
		salted = salted + dI
	}

	key := salted[0:32]
	iv := salted[32:48]

	return key, iv
}
