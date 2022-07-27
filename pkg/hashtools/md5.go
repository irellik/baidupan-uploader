package hashtools

import (
	"crypto/md5"
	"encoding/hex"
)

func HashMd5(strByte []byte) string {
	h := md5.New()
	h.Write(strByte)
	return hex.EncodeToString(h.Sum(nil))
}
