package experiment

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"tarea3/sa"
	"tarea3/tsp"
)


func Run(outPath string, id int64) error {

	rng := rand.New(rand.NewSource(id))

	r := rng.Float64()

	const TMin, TMax = 1, 50

	fracs := []float64{0.2*r, 0.2*r+0.2, 0.2*r+0.4, 0.2*r+0.6, 0.2*r+0.8}

	for i, f := range fracs {
		fracs[i] = TMin + f * (TMax - TMin)
	}

	file, errOS := os.Create(outPath)
	if errOS != nil {
		return errOS
	}

	defer file.Close()

	dataW := csv.NewWriter(file)
	dataW.Write([]string{"n", "nivel", "T0", "replica", "costo"})

	defer dataW.Flush()

	levelNames := []string{"axial_bajo", "bajo", "centro", "alto", "axial_alto"}

	for _, n := range []int{5, 10, 15} {

		in := tsp.NewInstance(n, 42)

		for c := range 5 {
			for d := range 10 {
				T := fracs[c]

				seed := int64(100*n + d)
				cost := sa.Run(in, T, seed, 100)

				dataW.Write([]string{fmt.Sprint(n), levelNames[c], fmt.Sprintf("%.6f", T), fmt.Sprint(d), fmt.Sprintf("%.6f", cost)})
			}
		}
	}

	return dataW.Error()
}
