package crypt

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Hash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Check(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}

func MD5Check(str, hash string) bool {
	// 计算明文 MD5
	sum := md5.Sum([]byte(str))
	calculated := hex.EncodeToString(sum[:])
	// 统一小写，避免大小写差异
	hash = strings.ToLower(hash)
	// 使用常量时间比较，防止时序攻击
	return subtle.ConstantTimeCompare([]byte(calculated), []byte(hash)) == 1
}
