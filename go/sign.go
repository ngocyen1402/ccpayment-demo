package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	PublicKey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCyPWejY7A+stkupI5Ow1aqlDgQ8g04gByyuyOiqw/wl8j8maerG1e7YKiF5qGOKr+Jw83HPdMFLCZDZebS63taPA2aIA+2x1CpIVfss5jSRQNsVzez9eDW7HTI+Nplx95BLl8OVE724hCgWFEjpwZ4GzORQMzmIXxxw67sdo9iuwIDAQAB\n-----END PUBLIC KEY-----\n"
	PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALI9Z6NjsD6y2S6kjk7DVqqUOBDyDTiAHLK7I6KrD/CXyPyZp6sbV7tgqIXmoY4qv4nDzcc90wUsJkNl5tLre1o8DZogD7bHUKkhV+yzmNJFA2xXN7P14NbsdMj42mXH3kEuXw5UTvbiEKBYUSOnBngbM5FAzOYhfHHDrux2j2K7AgMBAAECgYBXtQGfk/lxEN7wJcdlGJg3/hGMvR8mU1xL0uyZKiYA1R/wtMed2imUqd6jbTbIV17DMte6mECThgMaHTW1Smz6yrXYwPLmorkZmDxC4ggpvriH7sDgvBL++lOlLfRQqL7XLx72ZDaFWC0qFokKc5vviXBqWnTVMf/SQenSZGkgEQJBAN5z1x9Dyv2XyYwyJqXzEHWmvx7jjwqGQx6nFWnIVfeXQyJSSY7tqT6J4fGHe9eq5nbnqQo964RrR91Q+2iRGMkCQQDNHqjvgoT/skAXy80BP2Mt5W5pFjjeVlaCoaf006mTngkfB24ZmvxoxX5NfNBEGB/iS2KCsU5/h1ykpU3Lj+VjAkA9MwVl9pKr/cxXI5z6XsqSc5N0/gnmTVW94x3DAniUKysvEBBon/3F1M0yU6HAjaXl5Ine5XYb8h/NRXBFLlXxAkEAub1muqOU7bmqoiGxPMz6cWgNh+lQi7zgz5+06FT2fK6hkdB3mYYnxHP5wA8ixFaYIKGkzbXi4EZh1NG/VXKzAwJAFp+hcKz9oRO1LodExpdmATTd031g53X+3MMKG+PJREjAnC9wQL4RsmbzYP5NZ2dORIpNgRWawF2b1KJxWiiCsg==\n-----END PRIVATE KEY-----\n"
)

func RsaVerySignWithSha256(data, signData []byte, publicKey string) bool {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return false
	}

	hashed := sha256.Sum256(data)

	//sig, _ := base64.StdEncoding.DecodeString(string(signData))  base64 decode
	sig, _ := hex.DecodeString(string(signData)) // hex decode

	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], sig)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func RsaSignWithSha256(data []byte, keyBytes []byte) (string, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return "", errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err

	}
	str := hex.EncodeToString(signature)
	return str, nil
}

func RsaEncrypt(data, keyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKIXPublicKey Err:", err)
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		fmt.Println("EncryptPKCS1v15 Err:", err)
		return nil, err
	}
	return ciphertext, nil
}

func RsaDecrypt(ciphertext, keyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, err
	}
	return data, nil
}
