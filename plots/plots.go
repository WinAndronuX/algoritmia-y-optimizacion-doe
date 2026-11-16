package plots

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// QuadraticCurve: puntos observados + parabola ajustada + linea vertical en T*.
func QuadraticCurve(Tobs, Yobs []float64, b0, b1, b2, topt float64, valid bool, outPath string) error {

	p := plot.New()

	p.Title.Text = "Costo vs Temperatura (RSM)"
	p.X.Label.Text = "Temperatura"
	p.Y.Label.Text = "Costo"

	// (ejes se ajustan despues de agregar los datos)

	pts := make(plotter.XYs, len(Tobs))
	for i := range Tobs {
		pts[i].X = Tobs[i]
		pts[i].Y = Yobs[i]
	}

	sc, _ := plotter.NewScatter(pts)
	sc.GlyphStyle.Color = color.RGBA{R: 31, G: 119, B: 180, A: 255}
	p.Add(sc)

	if len(Tobs) > 0 {

		tmin, tmax := Tobs[0], Tobs[0]

		for _, t := range Tobs {

			if t < tmin {
				tmin = t
			}

			if t > tmax {
				tmax = t
			}
		}

		nC := 100
		cp := make(plotter.XYs, nC)

		for i := range nC {
			t := tmin + (tmax-tmin)*float64(i)/float64(nC-1)
			cp[i].X = t
			cp[i].Y = b0 + b1*t + b2*t*t
		}

		l, _ := plotter.NewLine(cp)
		l.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

		p.Add(l)
	}

	if valid {

		vl := plotter.XYs{{X: topt, Y: 0}, {X: topt, Y: 1e9}}
		lv, _ := plotter.NewLine(vl)

		lv.Color = color.RGBA{R: 0, G: 128, B: 0, A: 255}
		lv.Dashes = []vg.Length{vg.Points(4), vg.Points(2)}

		p.Add(lv)
	}

	// ajustar ejes al rango de los datos DESPUES de agregar (para que no se sobrescriba)
	if len(Tobs) > 0 {

		tmin, tmax := Tobs[0], Tobs[0]
		ymin, ymax := Yobs[0], Yobs[0]

		for i := range Tobs {
			if Tobs[i] < tmin { tmin = Tobs[i] }
			if Tobs[i] > tmax { tmax = Tobs[i] }
			if Yobs[i] < ymin { ymin = Yobs[i] }
			if Yobs[i] > ymax { ymax = Yobs[i] }
		}

		margin := (ymax - ymin) * 0.1

		if margin == 0 { margin = 1 }

		p.X.Min = tmin - (tmax-tmin)*0.05
		p.X.Max = tmax + (tmax-tmin)*0.05

		p.Y.Min = ymin - margin
		p.Y.Max = ymax + margin
	}

	return p.Save(7*vg.Inch, 4*vg.Inch, outPath)
}
