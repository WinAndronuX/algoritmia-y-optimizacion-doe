package tsp

import (
	"math"
	"math/rand"
)

type Point struct {
	X float64
	Y float64
}

type Instance struct {
	Points []Point
	N      int
}

func NewInstance(n int, seed int64) Instance {

	var instance = Instance{N: n}

	rng := rand.New(rand.NewSource(seed))

	for _ = range n {
		instance.Points = append(instance.Points, Point{X: rng.Float64() * 100, Y: rng.Float64() * 100})
	}

	return instance
}

func Dist(a, b Point) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

func TourCost(in Instance, tour []int) float64 {

	var cost float64

	for i := range in.N {
		cost += Dist(in.Points[tour[i]], in.Points[tour[(i+1)%in.N]])
	}

	return cost
}
