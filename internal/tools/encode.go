package tools

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// HashPassword 用于将密码哈希
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("哈希密码:%v失败", err)
	}
	return string(hashed), nil

}

// CheckPassword 检查密码
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),
		[]byte(password))
}

const SessionPrefix = "sess_"

// GetRandomToken 生成随机的字符串
func GetRandomToken(len int) string {
	r := make([]byte, len)
	io.ReadFull(rand.Reader, r)
	//再编码
	return base64.URLEncoding.EncodeToString(r)
}

// CreateSessionID 就是将session前缀和随机的Token组合
func CreateSessionID(token string) string {
	return SessionPrefix + token
}

func GetSessionIdByUserId(userId int) string {
	return fmt.Sprintf("sess_map_%d", userId)
}
