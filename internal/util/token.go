package util

import (
	"crypto/md5"
	"fmt"
	"time"
)

func GenerateToken(input string, salt string) string {
	data := []byte(fmt.Sprintf("%s%s%s", input, salt, time.Now().Format("2006_01_02_15_04_05")))
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", has)

	return md5Str
}
