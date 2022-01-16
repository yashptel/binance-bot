package utils

import "math"

func GetPercentageDiff(v1, v2 float64) float64 {
	return (math.Abs(v1-v2) / ((v1 + v2) / 2)) * 100
}
