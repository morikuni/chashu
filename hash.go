package chashu

import (
	"crypto/md5"
	"encoding/binary"
)

func defaultHash(key string) uint32 {
	r := md5.Sum([]byte(key))
	return binary.LittleEndian.Uint32(r[8:12])
}
