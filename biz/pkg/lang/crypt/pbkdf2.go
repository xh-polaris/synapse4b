package crypt

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 15 // 迭代次数
	keyLength  = 32 // 256 bits = 32 bytes
	saltLength = 32 // 32 bytes
)

// PBKDF2WithHmacSHA1 salt+hash加密
func PBKDF2WithHmacSHA1(str, salt string) (string, error) {
	var saltBytes []byte
	var err error

	if salt == "" {
		// 生成随机salt
		saltBytes = make([]byte, saltLength)
		if _, err = rand.Read(saltBytes); err != nil {
			return "", err
		}
	} else {
		// 使用传入的salt（用于迁移场景）
		if saltBytes, err = base64.StdEncoding.DecodeString(salt); err != nil {
			return "", err
		}
	}

	// 使用PBKDF2生成哈希
	hash := pbkdf2.Key([]byte(str), saltBytes, iterations, keyLength, sha1.New)

	// Base64编码
	base64Salt := base64.StdEncoding.EncodeToString(saltBytes)
	base64Hash := base64.StdEncoding.EncodeToString(hash)

	return base64Salt + ":" + base64Hash, nil
}

// PBKDF2WithHmacSHA1Check 用于校验salt+hash加密的密码
// str是明文密码，storedPwd是存储的`salt:密文串`格式的密码
func PBKDF2WithHmacSHA1Check(str, storedPwd string) bool {
	if str == "" || storedPwd == "" {
		return false
	}
	parts := strings.Split(storedPwd, ":")

	checkPwd, _ := PBKDF2WithHmacSHA1(str, parts[0])

	return checkPwd == storedPwd
}
