package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateID(prefix string) string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	//timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := time.Now().Format("20060102150405")
	orderID := node.Generate().Int64()

	uniqueID := fmt.Sprintf("%s%s%d", prefix, timestamp, orderID)
	return uniqueID
}
