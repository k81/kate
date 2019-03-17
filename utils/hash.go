package utils

import (
	"hash/crc32"
	"hash/fnv"
)

// CRC32 return the CRC32 hash of string
func CRC32(data string) uint32 {
	h := crc32.NewIEEE()
	// nolint:errcheck
	h.Write([]byte(data))
	return h.Sum32()
}

// FNV32a return the FNV32a hash of string
func FNV32a(data string) uint32 {
	h := fnv.New32a()
	// nolint:errcheck
	h.Write([]byte(data))
	return h.Sum32()
}
