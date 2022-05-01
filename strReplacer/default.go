package strReplacer

import (
	"math/rand"
	"time"
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
		if s, ok := letters[v]; ok { // 存在
			str += s
			continue
		}
		if v == "" {
			return ""
		}
		str += v // 自定义字符串
	}
	// 生成随机字符串
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b, strByt := make([]byte, n), []byte(str)
	for i := range b {
		b[i] = strByt[r.Intn(len(str))]
	}
	return string(b)
}
