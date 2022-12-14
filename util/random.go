package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	letters := RandomInt(6, 15)
	owner := RandomString(int(letters))
	return owner
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EURO", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandomMoney() int64 {
	return RandomInt(0, 100000)
}
