package tools

import "fmt"

// CreateSessionID 就是将session前缀和随机的Token组合
func CreateSessionID(token string) string {
	return SessionPrefix + token
}

func GetSessionIdByUserId(userId int) string {
	return fmt.Sprintf("sess_map_%d", userId)
}
