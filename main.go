package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"

	"tarea3/experiment"
	"tarea3/plots"
	"tarea3/stats"
)

func main() {
	_ = os.MkdirAll("data", 0755)
	_ = os.MkdirAll("figures", 0755)

	// 1) experimento
	if err := experiment.Run("data/results.csv", 2222332); err != nil {
		fmt.Println("error experimento:", err)
		return
	}
	fmt.Println("experimento listo")

	// 2) leer CSV
	f, err := os.Open("data/results.csv")
	if err != nil {
		fmt.Println("error abriendo csv:", err)
		return
	}

	defer f.Close()

	rd := csv.NewReader(f)

	rows, err := rd.ReadAll()
	if err != nil {
		fmt.Println("error leyendo csv:", err)
		return
	}

	// 3) agrupar por tamaño y nivel
	levelCosts := map[int]map[string][]float64{}
	levelTemp := map[int]map[string]float64{}

	var sizes []int
	sizeSet := map[int]bool{}

	for _, r := range rows[1:] {

		n, _ := strconv.Atoi(r[0])
		level := r[1]

		T, _ := strconv.ParseFloat(r[2], 64)
		cost, _ := strconv.ParseFloat(r[4], 64)

		if levelCosts[n] == nil {
			levelCosts[n] = map[string][]float64{}
			levelTemp[n] = map[string]float64{}
		}

		levelCosts[n][level] = append(levelCosts[n][level], cost)
		levelTemp[n][level] = T

		if !sizeSet[n] {
			sizeSet[n] = true
			sizes = append(sizes, n)
		}
	}

	sort.Ints(sizes)

	levels := []string{"axial_bajo", "bajo", "centro", "alto", "axial_alto"}

	// 4) RSM por tamaño + bootstrap + graficas
	fmt.Println()
	fmt.Println("=== RSM por tamaño ===")

	type summary struct {
		n       int
		bestLevel string
		bestMedian float64
		bestLo  float64
		bestHi  float64
		topt    float64
		valid   bool
		pb2     float64
		r2      float64
	}

	var summaries []summary

	for _, n := range sizes {

		var Tobs, Yobs []float64
		for _, lvl := range levels {
			T := levelTemp[n][lvl]
			for _, c := range levelCosts[n][lvl] {
				Tobs = append(Tobs, T)
				Yobs = append(Yobs, c)
			}
		}

		q := stats.FitQuadratic(Tobs, Yobs)
		fmt.Printf("n=%d: %s\n", n, q.String())

		bestMedian := math.Inf(1)
		bestLevel := ""

		var bestLo, bestHi float64

		for _, lvl := range levels {

			costs := levelCosts[n][lvl]
			med, lo, hi := stats.BootstrapMedian(costs, 10000, 42+int64(n))
			fmt.Printf("  %-12s  mediana=%.4f  IC=[%.4f, %.4f]\n", lvl, med, lo, hi)

			if med < bestMedian {
				bestMedian = med
				bestLevel = lvl
				bestLo = lo
				bestHi = hi
			}
		}

		fmt.Println()

		_ = plots.QuadraticCurve(Tobs, Yobs, q.B0, q.B1, q.B2, q.TOpt, q.Valid,
			fmt.Sprintf("figures/rsm_n%d.png", n))

		summaries = append(summaries, summary{
			n: n, bestLevel: bestLevel, bestMedian: bestMedian, bestLo: bestLo, bestHi: bestHi,
			topt: q.TOpt, valid: q.Valid, pb2: q.PB2, r2: q.R2,
		})
	}

	// 5) tabla final
	fmt.Println("=== Tabla resumen ===")
	fmt.Printf("%-6s %-14s %12s %22s %12s %8s %8s\n", "n", "Mejor nivel", "Mediana", "IC 95%", "T*", "R2", "valido")

	for _, r := range summaries {
		fmt.Printf("%-6d %-14s %12.4f [%8.4f, %8.4f] %12.4f %8.4f %8v\n",
			r.n, r.bestLevel, r.bestMedian, r.bestLo, r.bestHi, r.topt, r.r2, r.valid)
	}

	fmt.Println()
	fmt.Println("Analisis: revisa si T* y el mejor nivel cambian con n.")
}
