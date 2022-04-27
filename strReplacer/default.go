package strReplacer

import (
	"math/rand"
)

// 随机字符串
// 替换随机字符串, n为随机字符串的长度
//	字符串 "number"或"lower"或"upper"或"special"或自定义字符串
//	var letter []string
//	预设随机字符串
//	letters = map[string]string{
//		"number": "0123456789",
//		"lower":  "abcdefghijklmnopqrstuvwxyz",
//		"upper":  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
//	    "special": "!@#$%^&*"
//	}
func RandStr(letter []string, n int) string {
	// 预设随机字符串
	letters := map[string]string{
		"number":  "0123456789",
		"lower":   "abcdefghijklmnopqrstuvwxyz",
		"upper":   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"special": "!@#$%^&*",
	}
	// 处理letter
	var str string
	for _, v := range letter {
		if letters[v] != "" { // 从预设中取出
			str += letters[v]
		} else {
			if v != "" { // 自定义
				str += v
			} else {
				return "" // 没有字符串
			}
		}
	}
	// 生成随机字符串
	var strByt = []byte(str)
	b := make([]byte, n)
	for i := range b {
		b[i] = strByt[rand.Intn(len(str))]
	}
	return string(b)
}
