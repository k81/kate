package utils

import (
	"encoding/hex"

	"github.com/rogpeppe/fastuuid"
)

var uuid_generator = fastuuid.MustNewGenerator()

func FastUUID() [24]byte {
	return uuid_generator.Next()
}

func FastUUIDStr() string {
	b := uuid_generator.Next()
	return hex.EncodeToString(b[:])
}
