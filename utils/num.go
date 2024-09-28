package utils

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/iface"
	"math"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 向上取整float64
func CeilFloat64(x float64) int {
	return int(math.Ceil(RoundFloat(x, 2)))
}

func FnvNumber[T iface.Integer](x T) uint64 {
	var h = uint64(enums.Offset64)
	h *= enums.Prime64
	h ^= uint64(x)
	return h
}

// BinaryPow 返回2的n次方
func BinaryPow(n int) int {
	var ans = 1
	for i := 0; i < n; i++ {
		ans <<= 1
	}
	return ans
}

func IsSameSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
