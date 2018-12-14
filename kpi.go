package main
/***
Коэффициент приращения
 */

const (
	minKpi float32 = .2 // 200ms
	maxKpi float32 = 300 // 5 min
	startKpi float32 = 1 //start value
)

type Kpi struct {
	min float32
	max float32
	current float32
}

func fill() (Kpi) {
	return Kpi{current: startKpi, min: minKpi, max: maxKpi}
}

func (kpi *Kpi) stepDown() {
	if kpi.current /= 2;
		kpi.current < minKpi {
			kpi.current = minKpi
	}
}

func (kpi *Kpi) stepUp() {
	if kpi.current *= 2;
		kpi.current > maxKpi {
			kpi.current = maxKpi
	}
}

func (kpi *Kpi) forceDown() {
	if kpi.current /= 8;
		kpi.current < minKpi {
		kpi.current = minKpi
	}
}

func (kpi *Kpi) forceUp()  {
	if kpi.current *= 2 * 2 * 2;
		kpi.current > maxKpi {
		kpi.current = maxKpi
	}
}

