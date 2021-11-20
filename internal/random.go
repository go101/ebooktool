package internal

import (
	"crypto/rand"
	"encoding/hex"
	mathrand "math/rand"
	"time"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	rs := make([]byte, n*2)
	_, err := rand.Read(rs)
	if err != nil {
		mathrand.Read(rs)
	}
	return hex.EncodeToString(rs)
}
