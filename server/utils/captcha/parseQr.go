package captcha

import (
	"fmt"
	"github.com/tuotoo/qrcode"
	"io"
	"net/http"
	"os"
)

func ParseQrCodeImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("打开图片文件出错:", err)
		return "", err
	}
	defer file.Close()

	qr, err := qrcode.Decode(file)
	if err != nil {
		fmt.Println("解码二维码出错:", err)
		return "", err
	}

	fmt.Println("解析得到的文本:", qr.Content)
	return qr.Content, nil
}

func ParseRemoteQrCodeImage(imageURL string) (string, error) {
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("下载图片出错:", err)
		return "", err
	}
	defer response.Body.Close()

	file, err := os.Create("qrcode.png")
	if err != nil {
		fmt.Println("创建临时文件出错:", err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("保存图片文件出错:", err)
		return "", err
	}

	qr, err := qrcode.Decode(file)
	if err != nil {
		fmt.Println("解码二维码出错:", err)
		return "", err
	}

	fmt.Println("解析得到的文本:", qr.Content)
	return qr.Content, nil
}
