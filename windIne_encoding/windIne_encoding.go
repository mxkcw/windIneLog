//go:build !windows
// +build !windows

/*
Package windIne_encoding 编码转换"目前仅支持gbk 转utf-8"
*/
package windIne_encoding

// ConvertToUTF8UsedLocalENV 将本地获取到的字符串强制转为UTF-8。
func ConvertToUTF8UsedLocalENV(str string) (string, error) {
	return str, nil
}
