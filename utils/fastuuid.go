package utils

import (
	"encoding/hex"

	"github.com/rogpeppe/fastuuid"
)

var uuidGenerator = fastuuid.MustNewGenerator()

// FastUUID generate a new UUID as []byte
func FastUUID() [24]byte {
	return uuidGenerator.Next()
}

// FastUUIDStr generate a new UUID as string
func FastUUIDStr() string {
	b := uuidGenerator.Next()
	return hex.EncodeToString(b[:])
}
