package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

// HmacSHA256 return HMAC-SHA256 hash
func HmacSHA256(secret []byte, plain []byte) ([]byte, error) {
	return doHmac(sha256.New, secret, plain)
}

// HmacMD5 return HMAC-MD5 hash
func HmacMD5(secret []byte, plain []byte) ([]byte, error) {
	return doHmac(md5.New, secret, plain)
}

// HmacSHA1 return HMAC-SHA1 hash
func HmacSHA1(secret []byte, plain []byte) ([]byte, error) {
	return doHmac(sha1.New, secret, plain)
}

// SHA1 compute sha1 hash
func SHA1(data []byte) []byte {
	h := sha1.New()

	_, _ = h.Write(data)

	return h.Sum(nil)
}

func doHmac(h func() hash.Hash, secret []byte, plain []byte) ([]byte, error) {
	mac := hmac.New(h, secret)
	_, err := mac.Write(plain)
	if err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}
