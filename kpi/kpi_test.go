package kpi

import "testing"

func TestStepDown(t *testing.T) {

	kpi := NewKpi()

	if kpi.current != 1 {
		t.Error("Expected 1, got", kpi.current)
	}

	kpi.StepDown()
	if kpi.current != .5 {
		t.Error("Expected 0.5, got", kpi.current)
	}

	kpi.StepDown()
	if kpi.current != .25 {
		t.Error("Expected 0.25, got", kpi.current)
	}

	kpi.StepDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

	kpi.StepDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

}

func TestStepUp(t *testing.T) {

	kpi := NewKpi()

	if kpi.current != 1 {
		t.Error("Expected 1, got", kpi.current)
	}

	kpi.StepUp()
	if kpi.current != 2 {
		t.Error("Expected 2, got", kpi.current)
	}

	kpi.StepUp()
	if kpi.current != 4 {
		t.Error("Expected 4, got", kpi.current)
	}

	kpi.StepUp()
	if kpi.current != 8 {
		t.Error("Expected 8, got", kpi.current)
	}

	kpi.StepUp()
	if kpi.current != 16 {
		t.Error("Expected 16, got", kpi.current)
	}

	kpi.StepUp()
	kpi.StepUp()
	kpi.StepUp()
	kpi.StepUp()
	if kpi.current != 256 {
		t.Error("Expected 256, got", kpi.current)
	}

	kpi.StepUp()
	if kpi.current != maxKpi {
		t.Error("Expected 300, got", kpi.current)
	}
}

func TestStepUpDown(t *testing.T) {

	kpi := NewKpi()
	kpi.StepUp()   //2
	kpi.StepUp()   //4
	kpi.StepUp()   //8
	kpi.StepDown() //4

	if kpi.current != 4 {
		t.Error("Expected 4, got", kpi.current)
	}

	kpi.ForceUp() //32
	if kpi.current != 32 {
		t.Error("Expected 32, got", kpi.current)
	}

	kpi.StepDown() // 16
	if kpi.current != 16 {
		t.Error("Expected 16, got", kpi.current)
	}

	kpi.ForceDown() // 2
	if kpi.current != 2 {
		t.Error("Expected 2, got", kpi.current)
	}

	kpi.ForceDown()
	kpi.ForceDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

	kpi.ForceUp()
	kpi.ForceUp()
	kpi.ForceUp()
	kpi.ForceUp()
	kpi.ForceUp()

	if kpi.current != maxKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}
}
