package main

import (
	"common"
	"crypto/md5"
	"fmt"
)

// GenSig 生成签名
func GenSig(timestamp int64) string {
	h := md5.New()
	str := fmt.Sprintf("%s%d", common.MSDKKey, timestamp)
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
