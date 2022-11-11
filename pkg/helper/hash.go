package helper

import (
	"crypto/md5"
	"fmt"
)

func Md5(input string) string {
	has := md5.Sum([]byte(input))
	md5Str := fmt.Sprintf("%x", has)

	return md5Str
}
