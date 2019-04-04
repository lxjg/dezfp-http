package tools

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

// GeneratePassword 生成密码
func GeneratePassword(serialNumber string) string {
	randomStr := GetRandomString(10)
	tempStr := strings.Join([]string{randomStr, serialNumber}, "")
	h := md5.New()
	h.Write([]byte(tempStr))
	password := hex.EncodeToString(h.Sum(nil))
	return randomStr + base64.StdEncoding.EncodeToString([]byte(password))
}

// GetRandomString get random string
func GetRandomString(length uint) string {
	return getRandomKey("0123456789abcdefghijklmnopqrstuvwxyz", length)
}

// getRandomKey
func getRandomKey(str string, length uint) string {
	if length == 0 {
		return ""
	}

	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
