package similarity

import (
	"fmt"
	"math"

	"github.com/spf13/cast"
)

func CosineSimilarity[T float32 | float64 | int | int64 | int32](a []T, b []T) (cosine float64, err error) {
	const power = 2

	aLen := len(a)
	bLen := len(b)

	count := max(aLen, bLen)

	var (
		sum float64
		sa  float64
		sb  float64
	)

	for i := 0; i < count; i++ {
		if i < aLen {
			sa += math.Pow(cast.ToFloat64(a[i]), power)
		}
		if i < bLen {
			sb += math.Pow(cast.ToFloat64(b[i]), power)
		}
		if i < aLen && i < bLen {
			sum += cast.ToFloat64(a[i]) * cast.ToFloat64(b[i])
		}
	}

	if sa == 0 || sb == 0 {
		return 0, fmt.Errorf("vectors should not be null (all zeros)")
	}

	return sum / (math.Sqrt(sa) * math.Sqrt(sb)), nil
}

func PearsonCorrelation[T float32 | float64 | int | int64 | int32](a []T, b []T) (float64, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 && bLen == 0 {
		return 0, fmt.Errorf("input vectors are both empty")
	}
	count := min(aLen, bLen)
	if count == 0 {
		return 0, fmt.Errorf("no overlap between vectors")
	}

	var (
		sumA  float64
		sumB  float64
		sumA2 float64
		sumB2 float64
		sumAB float64
	)

	for i := 0; i < count; i++ {
		va := cast.ToFloat64(a[i])
		vb := cast.ToFloat64(b[i])
		sumA += va
		sumB += vb
		sumA2 += va * va
		sumB2 += vb * vb
		sumAB += va * vb
	}

	n := float64(count)
	numerator := sumAB - (sumA * sumB / n)
	denominator := math.Sqrt((sumA2 - (sumA*sumA)/n) * (sumB2 - (sumB*sumB)/n))
	if denominator == 0 {
		return 0, fmt.Errorf("cannot calculate correlation (division by zero)")
	}

	return numerator / denominator, nil
}

func HammingDistance(a, b string) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("lengths are not equal")
	}
	dist := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			dist++
		}
	}
	return dist, nil
}
func ManhattanDistance[T int | int64 | float32 | float64](a, b []T) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have the same length")
	}
	sum := 0.0
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		if diff < 0 {
			diff = -diff
		}
		sum += diff
	}
	return sum, nil
}
