package captcha

import (
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"go.uber.org/zap"
	"os"
	"strings"
)

func ParseQrCodeImageFromBase64(base64String string) (string, error) {
	// 通过分隔符检测图片格式
	var format string
	if strings.Contains(base64String, "data:image/jpeg;base64") {
		format = "jpeg"
	} else if strings.Contains(base64String, "data:image/png;base64") {
		format = "png"
	} else if strings.Contains(base64String, "data:image/gif;base64") {
		format = "gif"
	} else if strings.Contains(base64String, "data:image/bmp;base64") {
		format = "bmp"
		// 其他图片格式判断可以继续添加
	} else {
		return "", fmt.Errorf("Unsupported image format")
	}

	global.GVA_LOG.Info("当前图片格式：", zap.Any("format", format))

	// 分离base64编码部分
	imageB64 := strings.Split(base64String, ";base64,")[1]

	imageBytes, err := base64.StdEncoding.DecodeString(imageB64)
	if err != nil {
		fmt.Println("解码 Base64 图片出错:", err)
		return "", err
	}

	file, err := os.Create("temp_image.jpg")
	if err != nil {
		fmt.Println("创建临时图片文件出错:", err)
		return "", err
	}
	defer file.Close()

	_, err = file.Write(imageBytes)
	if err != nil {
		fmt.Println("写入临时图片文件出错:", err)
		return "", err
	}

	return ParseQrCodeImage("temp_image.jpg")
}

func ParseQrCodeImage(imagePath string) (string, error) {
	//file, err := os.Open(imagePath)
	//if err != nil {
	//	fmt.Println("打开图片文件出错:", err)
	//	return "", err
	//}
	//defer file.Close()

	img, err := imaging.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
	}

	// 识别二维码
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		return "", fmt.Errorf("无法解析改图二维码，请核查")
	}
	fmt.Println(result)

	return result.String(), nil
}

/*func ParseRemoteQrCodeImage(imageURL string) (string, error) {
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
}*/
