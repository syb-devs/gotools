package rand

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	Seed(time.Now().UTC().UnixNano())
}

func Seed(seed int64) {
	rand.Seed(seed)
}

func AlphaNum(size int) string {
	slice := make([]byte, size)
	alphanum := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	numChars := len(alphanum)

	for i := 0; i < size; i++ {
		slice[i] = alphanum[rand.Intn(numChars)]
	}
	return fmt.Sprintf("%s", slice)
}
