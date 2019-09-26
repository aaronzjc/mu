package tool

import (
	"crawler/internal/util/config"
	"crypto/md5"
	"fmt"
	"time"
)

func GenerateToken(input string) string {
	data := []byte(fmt.Sprintf("%s%s%s", input, config.NewConfig().Salt, time.Now().Format("2006_01_02_15_04_05")))
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", has)

	return md5Str
}

func ArrSearch(ele string, arr []string) int {
	res := -1
	for idx, val := range arr {
		if val == ele {
			res = idx
			break
		}
	}

	return res
}