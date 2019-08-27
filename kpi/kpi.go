//Коэффициент приращения
package kpi

const (
	minKpi   float32 = .2  // 200ms
	maxKpi   float32 = 300 // 5 min
	startKpi float32 = 1   //start value
)

type kpi struct {
	min     float32
	max     float32
	current float32
}

//NewKpi создает новый коэффициент приращения
func NewKpi() *kpi {
	return &kpi{
		min:     minKpi,
		max:     maxKpi,
		current: startKpi,
	}
}

func (k *kpi) GetCurrent() float32 {
	return k.current
}

//StepDown уменьшить коэффициент
func (k *kpi) StepDown() {
	if k.current /= 2; k.current < minKpi {
		k.current = minKpi
	}
}

// StepUp увеличить коэффициент
func (k *kpi) StepUp() {
	if k.current *= 2; k.current > maxKpi {
		k.current = maxKpi
	}
}

func (k *kpi) ForceDown() {
	if k.current /= 8; k.current < minKpi {
		k.current = minKpi
	}
}

func (k *kpi) ForceUp() {
	if k.current *= 2 * 2 * 2; k.current > maxKpi {
		k.current = maxKpi
	}
}
