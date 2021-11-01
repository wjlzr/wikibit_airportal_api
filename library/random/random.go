package random

import (
	"math/rand"
	"time"
)

var letters = []byte("abcdefghjkmnpqrstuvwxyz123456789")

func init() {
	rand.Seed(time.Now().Unix())
}

// Rand 随机字符串，包含 1~9 和 a~z - [i,l,o]
func Rand(n int) string {
	if n <= 0 {
		return "qc902ts5100twybq3y6r3x9lf61k06fb"
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return "qc902ts5100twybq3y6r3x9lf61k06fb"
	}
	for i, x := range b {
		arc = x & 31
		b[i] = letters[arc]
	}
	return string(b)
}
