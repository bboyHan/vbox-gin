package captcha

import (
	"errors"
	"fmt"
	"github.com/pquerna/otp/totp"
	"net/url"
)

func AuthQrCode(userName string) (string, error) {
	// 生成新的密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "VBOX",
		AccountName: userName,
	})
	if err != nil {
		fmt.Println("生成密钥出错:", err)
		return "", err
	}

	// 生成 TOTP URL
	_url := key.URL()
	//secret := key.Secret()

	// 生成二维码图片并保存到文件
	//err = qrcode.WriteFile(_url, qrcode.Medium, 256, "qrcode.png")
	//if err != nil {
	//	fmt.Println("生成二维码出错:", err)
	//	return "", err
	//}

	//fmt.Println("请使用双因子认证器扫描以下二维码进行认证:")
	fmt.Println(_url)
	//fmt.Println(secret)
	return _url, err
}

func ValidateCode(secret string, code string) bool {
	valid := totp.Validate(code, secret)
	if valid {
		fmt.Println("认证成功!")
	} else {
		fmt.Println("认证失败!")
	}
	return valid
}

func GetSecret(urlString string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("认证失败!")
		return "", errors.New("原认证码错误")
	}

	// 获取参数
	queryParams := parsedURL.Query()

	// 获取secret参数值
	secret := queryParams.Get("secret")

	fmt.Println("Secret:", secret)
	return secret, nil
}
