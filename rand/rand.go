package rand

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	SeedNow()
}

// Seed uses the provided seed value to initialize the default Source to a deterministic state.
func Seed(seed int64) {
	rand.Seed(seed)
}

// SeedNow uses the current timestamp as the seed
func SeedNow() {
	Seed(time.Now().UTC().UnixNano())
}

// AlphaNum generates a random alphanumeric string with the given length
func AlphaNum(size int) string {
	slice := make([]byte, size)
	alphanum := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	numChars := len(alphanum)

	for i := 0; i < size; i++ {
		slice[i] = alphanum[rand.Intn(numChars)]
	}
	return fmt.Sprintf("%s", slice)
}
