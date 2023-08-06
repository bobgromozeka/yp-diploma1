package functions

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func CheckLuhn(payload string) bool {
	reversed := Reverse(strings.Split(payload, ""))

	if len(reversed) < 1 {
		return false
	}

	sum := 0
	for i, strNum := range reversed {
		intNum, err := strconv.Atoi(strNum)
		if err != nil {
			return false
		}
		if i%2 != 0 {
			intNum *= 2
			if intNum > 9 {
				intNum -= 9
			}
		}
		sum += intNum
	}
	return sum%10 == 0
}

func Reverse[T constraints.Ordered](s []T) []T {
	reversed := make([]T, len(s))

	for i := 0; i < len(s); i++ {
		reversed[len(s)-i-1] = s[i]
	}

	return reversed
}
