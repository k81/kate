package utils

import (
	"hash/crc32"
	"hash/fnv"
)

func CRC32(data string) uint32 {
	h := crc32.NewIEEE()
	_, _ = h.Write([]byte(data))
	return h.Sum32()
}

func FNV32a(data string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(data))
	return h.Sum32()
}
