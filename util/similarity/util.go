package similarity

import (
	"errors"
	"fmt"
	"math"

	"github.com/karosown/katool-go/container/stream"
	"github.com/spf13/cast"
)

func CosineSimilarity[T ~float32 | ~float64 | ~int | ~int64 | ~int32](a []T, b []T) (cosine float64, err error) {
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

func PearsonCorrelation[T ~float32 | ~float64 | ~int | ~int64 | ~int32](a []T, b []T) (float64, error) {
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

func HammingDistance[T ~float32 | ~float64 | ~int | ~int64 | ~int32](a, b []T) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("lengths are not equal")
	}
	dist := 0
	for i := range a {
		if a[i] != b[i] {
			dist++
		}
	}
	return dist, nil
}
func HammingDistanceStr(a, b string) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("lengths are not equal")
	}
	dist := 0
	for i := range a {
		if a[i] != b[i] {
			dist++
		}
	}
	return dist, nil
}
func ManhattanDistance[T ~float32 | ~float64 | ~int | ~int64 | ~int32](a, b []T) (float64, error) {
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

type SimilarityFunc[T ~float32 | ~float64 | ~int | ~int64 | ~int32] func(a, b []T) (float64, error)

type Vector[T ~float32 | ~float64 | ~int | ~int64 | ~int32] interface {
	GetVector() []T
}

func TopK[T ~float32 | ~float64 | ~int | ~int64 | ~int32, V []T](k int, a V, b []V, similarityFunc SimilarityFunc[T]) ([]V, error) {
	if len(b) <= k {
		return b, nil
	}
	type KV struct {
		Key   V
		Value float64
	}
	s := stream.ToStream(&b).Map(func(i V) any {
		f, err := similarityFunc(a, i)
		if err != nil {
			return err
		}
		return KV{
			Key:   i,
			Value: f,
		}
	})
	s2 := stream.Cast[error](s.Filter(func(i any) bool {
		err, ok := i.(error)
		return ok && err != nil
	}))
	if s2.Count() > 0 {
		return nil, errors.Join(s2.ToList()...)
	}
	return stream.Cast[V](stream.Cast[KV](s.Filter(func(i any) bool {
		_, ok := i.(KV)
		return ok
	})).Sort(func(a, b KV) bool {
		return a.Value > b.Value
	}).Sub(0, k).Map(func(i KV) any {
		return i.Key
	})).ToList(), nil
}
func TopKByVector[R ~float32 | ~float64 | ~int | ~int64 | ~int32, T Vector[R], V []T](k int, a T, b V, similarityFunc SimilarityFunc[R]) (V, error) {
	if len(b) <= k {
		return b, nil
	}

	// 1. 流式获取全部向量和index映射
	rList := stream.Cast[[]R](stream.ToStream(&b).Map(func(t T) any { return t.GetVector() })).ToList()

	// 2. Map向量到对象T的映射，防止同向量多对象丢失
	indexMap := map[string]T{}
	for i, v := range rList {
		key := fmt.Sprint(v)             // 以向量内容做唯一key
		if _, ok := indexMap[key]; !ok { // 只记录第一个
			indexMap[key] = b[i]
		}
	}

	// 3. 用现有TopK，得TopK个向量
	topVecs, err := TopK[R, []R](k, a.GetVector(), rList, similarityFunc)
	if err != nil {
		return nil, err
	}

	// 4. 恢复回原始对象
	result := stream.Cast[T](stream.ToStream(&topVecs).Map(func(v []R) any {
		return indexMap[fmt.Sprint(v)]
	})).ToList()

	return result, nil
}
