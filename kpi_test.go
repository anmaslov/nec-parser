package main

import "testing"

func TestStepDown(t *testing.T) {

	kpi := fill()

	if kpi.current != 1 {
		t.Error("Expected 1, got", kpi.current)
	}

	kpi.stepDown()
	if kpi.current != .5 {
		t.Error("Expected 0.5, got", kpi.current)
	}

	kpi.stepDown()
	if kpi.current != .25 {
		t.Error("Expected 0.25, got", kpi.current)
	}

	kpi.stepDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

	kpi.stepDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

}

func TestStepUp(t *testing.T) {

	kpi := fill()

	if kpi.current != 1 {
		t.Error("Expected 1, got", kpi.current)
	}

	kpi.stepUp()
	if kpi.current != 2 {
		t.Error("Expected 2, got", kpi.current)
	}

	kpi.stepUp()
	if kpi.current != 4 {
		t.Error("Expected 4, got", kpi.current)
	}

	kpi.stepUp()
	if kpi.current != 8 {
		t.Error("Expected 8, got", kpi.current)
	}

	kpi.stepUp()
	if kpi.current != 16 {
		t.Error("Expected 16, got", kpi.current)
	}

	kpi.stepUp()
	kpi.stepUp()
	kpi.stepUp()
	kpi.stepUp()
	if kpi.current != 256 {
		t.Error("Expected 256, got", kpi.current)
	}

	kpi.stepUp()
	if kpi.current != maxKpi {
		t.Error("Expected 300, got", kpi.current)
	}
}

func TestStepUpDown(t *testing.T) {

	kpi := fill()
	kpi.stepUp() //2
	kpi.stepUp() //4
	kpi.stepUp() //8
	kpi.stepDown() //4

	if kpi.current != 4 {
		t.Error("Expected 4, got", kpi.current)
	}

	kpi.forceUp() //32
	if kpi.current != 32 {
		t.Error("Expected 32, got", kpi.current)
	}

	kpi.stepDown() // 16
	if kpi.current != 16 {
		t.Error("Expected 16, got", kpi.current)
	}

	kpi.forceDown() // 2
	if kpi.current != 2 {
		t.Error("Expected 2, got", kpi.current)
	}

	kpi.forceDown()
	kpi.forceDown()
	if kpi.current != minKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}

	kpi.forceUp()
	kpi.forceUp()
	kpi.forceUp()
	kpi.forceUp()
	kpi.forceUp()

	if kpi.current != maxKpi {
		t.Error("Expected 0.2, got", kpi.current)
	}
}