package main

//import (
//	"fmt"
//	"html"
//	"regexp"
//	"strconv"
//	"strings"
//)
//
//func decodeNonStandardUnicode(input string) string {
//	re := regexp.MustCompile(`%u([0-9a-fA-F]{4})`)
//	decoded := re.ReplaceAllStringFunc(input, func(s string) string {
//		unicodeStr := s[2:]
//		unicodeInt, err := strconv.ParseInt(unicodeStr, 16, 32)
//		if err != nil {
//			return s
//		}
//		return string(rune(unicodeInt))
//	})
//
//	return decoded
//}
//
//func main() {
//	hb := `<html><head><title>Object moved</title></head><body>
//<h2>Object moved to <a href="/Charge/UCardChargeInfo.aspx?ts=638464475251657885&amp;current=1&amp;bn=C240319614957121&amp;sn=8be375e50c0d760bbd4597798b22f6de&amp;m=%u5b98%u65b9%u5145%u503c%u51fa%u73b0%u672a%u77e5%u60c5%u51b5%uff0c%u539f%u56e0%uff1a%u521b%u5efa%u5355%u636e%u6210%u529f">here</a>.</h2>
//</body></html>
//`
//
//	// 解码HTML实体字符
//	decodedHTML := html.UnescapeString(hb)
//
//	// 替换 &amp; 为 &
//	decodedHTML = strings.ReplaceAll(decodedHTML, "&amp;", "&")
//
//	// 解码非标准Unicode编码
//	decodedHTML = decodeNonStandardUnicode(decodedHTML)
//
//	fmt.Println(decodedHTML)
//}
