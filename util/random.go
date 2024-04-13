package util

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand
const alphapbet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	s := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s)
}

func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphapbet)

	for i := 0; i < n; i++ {
		c := alphapbet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}


// RandomOwner returns a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "MYR"}
	n := len(currencies)
	return currencies[r.Intn(n)]
}

