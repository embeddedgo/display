package nightstrike

import (
	"math"
	"math/rand"
)

func rsin(r, th int) int {
	return int(float64(r) * math.Sin((2*math.Pi)*float64(th)/65536.0))
}

func rcos(r, th int) int {
	return rsin(r, th+0x4000)
}

func atan2(y, x int) int {
	return int(math.Atan2(float64(y), float64(x))*0x8000/math.Pi) - 0x4000
}

func abs(v int) int {
	if v < 0 {
		v = -v
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func prob(num, denom int) bool {
	return int(rand.Int63())&0xFFFF < num<<16/denom
}
