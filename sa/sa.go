package sa

import (
	"math"
	"math/rand"
	"tarea3/tsp"
)


func Run(in tsp.Instance, T0 float64, seed int64, iters int) float64 {

	tour := make([]int, in.N)
	for i := range in.N {
		tour[i] = i
	}

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(in.N, func(i, j int) {tour[i], tour[j] = tour[j], tour[i]})

	best := tsp.TourCost(in, tour)
	curr := best
	T := T0

	for _ = range iters {

		var (
			i = rng.Intn(in.N - 1)
			j = i + 1 + rng.Intn(in.N - i - 1)
			nb = make([]int, in.N)
		)

		copy(nb, tour)

		for n, m := i+1, j; n < m; n, m = n+1, m-1 {
			nb[n], nb[m] = nb[m], nb[n]
		}

		nbCost := tsp.TourCost(in, nb)

		delta := nbCost - curr

		if delta < 0 || (T > 0 && rng.Float64() < math.Exp(-delta/T)) {
			tour = nb
			curr = nbCost

			if curr < best {
				best = curr
			}
		}

		T *= 0.95
		if T < 1e-6 {
			T = 1e-6
		}
	}

	return best
}
