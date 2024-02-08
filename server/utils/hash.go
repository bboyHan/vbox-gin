package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

	uniqueID := fmt.Sprintf("%s%s%d", prefix, timestamp, orderID)[:18]
	return uniqueID
}

func Prefix(prefix string, target string) string {
	uniqueID := fmt.Sprintf("%s%s", prefix, target)
	return uniqueID
}

// RandomElement 从给定的字符串切片中随机选择一个元素并返回。
// 如果切片为空，将返回空字符串。
func RandomElement(elements []string) string {
	// 检查切片是否为空
	if len(elements) == 0 {
		return ""
	}

	// 初始化随机数生成器
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// 生成随机索引并返回对应的元素
	randomIndex := rng.Intn(len(elements))
	return elements[randomIndex]
}
