package utils

import (
	"crypto/md5"
	"encoding/hex"
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

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func GenerateID(prefix string) string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	//timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := time.Now().Format("20060102150405")
	orderID := node.Generate().Int64()

	uniqueID := fmt.Sprintf("%s%s%d", prefix, timestamp, orderID)[:24]
	return uniqueID
}
