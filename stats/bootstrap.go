package stats

import (
	"math/rand"
	"sort"
)

func median(xs []float64) float64 {

	cp := append([]float64(nil), xs...)
	sort.Float64s(cp)

	m := len(cp) / 2

	if len(cp)%2 == 0 {
		return (cp[m-1] + cp[m]) / 2
	}

	return cp[m]
}

// BootstrapMedian devuelve la mediana observada y el IC 95% por percentiles.
func BootstrapMedian(data []float64, B int, seed int64) (float64, float64, float64) {

	rng := rand.New(rand.NewSource(seed))
	meds := make([]float64, B)
	n := len(data)

	for b := range B {

		samp := make([]float64, n)

		for i := range n {
			samp[i] = data[rng.Intn(n)]
		}

		meds[b] = median(samp)
	}

	sort.Float64s(meds)

	lo := meds[int(0.025*float64(B))]
	hi := meds[int(0.975*float64(B))]

	return median(data), lo, hi
}
