package util

import (
	"math/rand"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func convertToStringMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range m {
		strKey := key.(string) // 假设 key 一定是字符串类型
		switch value := value.(type) {
		case map[interface{}]interface{}:
			result[strKey] = convertToStringMap(value) // 递归转换嵌套的 map
		default:
			result[strKey] = value
		}
	}
	return result
}
