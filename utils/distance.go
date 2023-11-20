package utils

import "math"

func CosineDistance(u []float64, v []float64) float64 {
	count := 0
	length_a := len(u)
	length_b := len(v)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(v[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(u[k], 2)
			continue
		}
		sumA += u[k] * v[k]
		s1 += math.Pow(u[k], 2)
		s2 += math.Pow(v[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2))
}
