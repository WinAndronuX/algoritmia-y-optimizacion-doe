# Tarea 3 – Parameter Tuning with Design of Experiments (DOE)

Adjusting the initial temperature (T₀) of Simulated Annealing (2-opt TSP) using Response Surface Methodology (RSM).

## Requirements

- Go 1.22+

## Run

```bash
go mod tidy
go run .
```

This generates:

- `data/results.csv` — 150 runs (3 sizes × 5 levels × 10 replicates)
- `figures/rsm_n5.png`, `rsm_n10.png`, `rsm_n15.png` — fitted curves
- Console output with RSM model, bootstrap medians/IC, and summary table

## Project Structure

```
Tarea3/
├── go.mod
├── main.go              # orchestrator: experiment + RSM + bootstrap + plots
├── tsp/tsp.go           # TSP instance + tour cost
├── sa/sa.go             # SA 2-opt (cooling 0.95)
├── experiment/run.go    # experimental design + CSV export
├── stats/bootstrap.go   # bootstrap median + 95% CI
├── stats/rsm.go         # quadratic fit + T*
├── plots/plots.go       # RSM plots
├── data/results.csv     # generated
├── figures/             # generated
├── REPORT.md
└── README.md
```

## Experiment Parameters

| Parameter | Value |
|---|---|
| Student ID (seed) | 2222332 |
| `r` (fixed) | 0.03956364384543931 |
| Temperature range | [1, 50] |
| Levels (5) | axial_low, low, center, high, axial_high |
| Replicates | 10 per design point |
| SA iterations | 100 |
| Cooling rate | 0.95 |
| Bootstrap B | 10,000 |

## Method

1. Generate `r` once from the student ID.
2. Define 5 temperature levels as fractions of [1, 50].
3. Run 10 replicates per design point with distinct seeds.
4. Fit a quadratic model `Cost = β₀ + β₁·T + β₂·T²` per problem size.
5. Verify β₂ > 0 for a valid minimum T* = −β₁/(2β₂).
6. Compute bootstrap median and 95% CI per design point.
7. Plot observed points + fitted curve + T* (when valid).

## Results Summary

| n | Best level | Median | 95% CI | T* | β₂ > 0 |
|---|---|---|---|---|---|
| 5 | axial_low | 193.43 | [193.43, 193.43] | — | No |
| 10 | axial_low | 290.99 | [289.77, 326.37] | — | No |
| 15 | axial_low | 426.34 | [410.43, 454.81] | 15.46 | Yes |

See `REPORT.md` for full analysis.
