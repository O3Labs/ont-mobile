package ontmobile

import "math"

func RoundFixed(val float64, decimals int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(decimals))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
