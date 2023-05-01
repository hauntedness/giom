package eth

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestDivideBigNumber(t *testing.T) {
	// for each random big number
	for i := 0; i < 10000; i++ {
		value := _rand_number()
		num1 := decimal.RequireFromString(value)
		num2 := WeiToEther(num1)
		s2 := strings.ReplaceAll(num2.String(), ".", "")
		if strings.Trim(s2, "0") != strings.Trim(num1.String(), "0") {
			t.Errorf("unexpected value, before division: %s, after division: %s", num1, num2)
		}
	}
}

func _rand_number() string {
	return fmt.Sprintf("%d%d%d", rand.Int63(), rand.Int63(), rand.Int63())
}

