package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	letters := randomInt(6, 15)
	owner := randomString(int(letters))
	return owner
}

func RandomCurrency() string {
	currencies := []string{"USD", "EURO", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandomMoney() int64 {
	return randomInt(0, 100000)
}
