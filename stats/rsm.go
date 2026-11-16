package stats

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

type QuadFit struct {
	B0, B1, B2 float64 // Costo = B0 + B1*T + B2*T^2
	PB2        float64 // p-value del termino cuadratico
	R2         float64 // coeficiente de determinacion
	TOpt       float64 // T* = -B1/(2*B2)
	Valid      bool    // B2 > 0
}

// FitQuadratic ajusta por minimos cuadrados y devuelve el modelo.
func FitQuadratic(T, Y []float64) QuadFit {

	n := len(T)
	X := mat.NewDense(n, 3, nil)

	for i := range n {
		X.Set(i, 0, 1)
		X.Set(i, 1, T[i])
		X.Set(i, 2, T[i]*T[i])
	}

	y := mat.NewVecDense(n, Y)

	var XtX mat.Dense
	XtX.Mul(X.T(), X)

	var XtY mat.VecDense
	XtY.MulVec(X.T(), y)

	var beta mat.VecDense
	if err := beta.SolveVec(&XtX, &XtY); err != nil {
		fmt.Println("error solve:", err)
		return QuadFit{}
	}

	b0, b1, b2 := beta.AtVec(0), beta.AtVec(1), beta.AtVec(2)

	var sse float64
	var yMean float64

	for i := range n {
		yMean += Y[i]
	}

	yMean /= float64(n)

	var sst float64

	for i := range n {

		pred := b0 + b1*T[i] + b2*T[i]*T[i]
		r := Y[i] - pred

		sse += r * r
		d := Y[i] - yMean

		sst += d * d
	}

	r2 := 1 - sse/sst
	p := 3
	sigma2 := sse / float64(n-p)

	var XtXinv mat.Dense
	if err := XtXinv.Inverse(&XtX); err != nil {
		return QuadFit{B0: b0, B1: b1, B2: b2}
	}

	seB2 := math.Sqrt(sigma2 * XtXinv.At(2, 2))
	tStat := b2 / seB2

	df := float64(n - p)
	pB2 := 2 * (1 - distuv.StudentsT{Nu: df}.CDF(math.Abs(tStat)))

	q := QuadFit{B0: b0, B1: b1, B2: b2, PB2: pB2, R2: r2}

	if b2 > 0 {
		q.TOpt = -b1 / (2 * b2)
		q.Valid = true
	}

	return q
}

func (q QuadFit) String() string {
	return fmt.Sprintf("Costo = %.4f + %.4f*T + %.4f*T^2  (p beta2=%.4g, R2=%.4f)  T*=%.4f  valido=%v",
		q.B0, q.B1, q.B2, q.PB2, q.R2, q.TOpt, q.Valid)
}
