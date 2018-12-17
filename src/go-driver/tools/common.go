package tools

import (
	"math"
)

/*
	计算除法后向商户取整
*/
func CountDividCeil(total int, deter int) int {

	time := math.Ceil(float64(total) / float64(deter))

	return int(time)
}
