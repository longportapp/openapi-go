package http

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func md5Hex(s string) []byte {
	data := md5.Sum([]byte(s))
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data[:])
	return dst
}

// EncryptPassword use to encrypt password
// Use aes with cbc mode to encrypt password
// 1. use md5(timestamp) as secret key for md5 have fixed length
// 2. use key[:16] is cbc iv
// 3. md5(raw_password) as plain content
// 4. padding plain content
// 5. encrypt

func EncryptPassword(raw, timestamp string) (encrypted string, err error) {
	if raw == "" {
		err = errors.New("empty password")
		return
	}

	key := md5Hex(timestamp)
	pwd := md5Hex(raw)
	iv := []byte(key)[:16]

	e, err := cbcEncryptWithIV(key, iv, pwd)
	if err != nil {
		return "", err
	}

	ret := base64.StdEncoding.EncodeToString(e)

	return ret, nil
}

// CBCEncryptWithIV use cbc cipher to encrypt data, but is is fixed size byte array
func cbcEncryptWithIV(key, iv, plain []byte) (encrypted []byte, err error) {
	var (
		block cipher.Block
	)

	padded := pkcs7Padding(plain, aes.BlockSize)

	encrypted = make([]byte, len(padded))

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, padded)

	return
}

func pkcs7Padding(b []byte, blocksize int) []byte {
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb
}
