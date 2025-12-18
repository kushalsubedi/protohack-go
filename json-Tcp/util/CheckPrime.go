package util

import "math"

func CheckPrime(content float64) bool {
	n := math.Floor(content)
	g := int(n)
	if n != content {
		return false
	}
	if g < 2 {
		return false
	}

	for i := 2; i*i <= g; i++ {
		if g%i == 0 {
			return false
		}
	}
	return true
}
