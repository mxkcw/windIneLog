package windIne_string

/*
Package windIne_string 字符串处理
*/
import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"regexp"
	"strings"
	"unsafe"
)

var simpleBytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

// RandomString 获取随机字符串
// n 多少位 比如获取10位随机数即传10
func RandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// WindIneCheckMobile 判断是否为手机号
func WindIneCheckMobile(phoneNumber string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phoneNumber)
}

// WindIneStringSliceContain 判断是否包含元素
func WindIneStringSliceContain(slice []string, s string) bool {
	for _, s2 := range slice {
		if s == s2 {
			return true
		}
	}
	return false
}

// WindIneValidHostnamePort 检测IP+端口字符串是否正确
func WindIneValidHostnamePort(s string) bool {
	sp := strings.Split(s, ":")
	if len(sp) != 2 {
		return false
	}
	if sp[0] == "" || sp[1] == "" {
		return false
	}
	return true
}

func WindIneRecodingString(src string, srcCode string, toCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(toCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// WindIneUTF8String2GBKString UTF-8转GBK
func WindIneUTF8String2GBKString(src string) string {
	return WindIneRecodingString(src, "GBK", "UTF-8")
}

// WindIneGBKString2UTF8String	GBK转UTF-8
func WindIneGBKString2UTF8String(src string) string {
	return WindIneRecodingString(src, "GBK", "UTF-8")
}

// WindIneBytes2String byte转string
func WindIneBytes2String(BytesData []byte) string {
	return *(*string)(unsafe.Pointer(&BytesData))
}

// WindIneString2Bytes string转bytes
func WindIneString2Bytes(strData string) []byte {
	return *(*[]byte)(unsafe.Pointer(&strData))
}

// WindIneStruct2JsonString struct转json
func WindIneStruct2JsonString(value interface{}) (jsonString string) {
	var cuValue, _ = json.Marshal(value)
	jsonString = string(cuValue)
	return jsonString
}

// DelStringEndNewlines 删除字符串结尾的 \n or \r\n
func DelStringEndNewlines(s *string) {
	b := []byte(*s)
	b = bytes.TrimSuffix(b, []byte("\r\n"))
	b = bytes.TrimSuffix(b, []byte("\n"))
	*s = string(b)
}

// StringCoverBool str value is true or false  covert to golang bool type
func StringCoverBool(str string) bool {

	return str == "true"
}
