package utils

import (
	"math/rand"
	"time"
)

func GenerateUID(length int, salt string) string {
	charSet := "0123456789"
	uid := make([]byte, length+len(salt))

	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串部分
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		uid[i] = charSet[randomIndex]
	}

	// 添加salt部分
	for i := 0; i < len(salt); i++ {
		uid[length+i] = salt[i]
	}

	return string(uid)
}

func GenerateContent(length int) string {

	charSet := "abcdefghijklmnopqrstuvwxyz"
	content := make([]byte, length)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		content[i] = charSet[randomIndex]
	}

	return string(content)
}
