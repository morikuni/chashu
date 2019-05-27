package chashu

import (
	"hash/fnv"
)

func fnvHash(key string) uint64 {
	h := fnv.New64()
	h.Write([]byte(key))
	// some salt seems to be needed for good distribution for short key
	h.Write([]byte("salt"))
	return h.Sum64()
}
