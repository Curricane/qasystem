package util

import (
	"crypto/md5"
	"fmt"
)

func Md5(data [] byte) (ret string) {
	md5Sum := md5.Sum(data)
	ret = fmt.Sprintf("%x", md5Sum)
	return
}
