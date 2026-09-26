package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString lenght n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner account
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney

func RandomMoney() int64 {
	return int64(rand.Intn(1000) + 100)

}

//RandomCurrency

func RandomCurrency() string {
	currencies := []string{
		"USD",
		"EUR",
		"THB",
	}
	return currencies[rand.Intn(len(currencies))]
}
