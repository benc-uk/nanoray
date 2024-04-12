package raytrace

import (
	"crypto/sha256"
	"fmt"
)

func GenerateID(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	id := fmt.Sprintf("%x", hash.Sum(nil))
	return id[:6]
}
