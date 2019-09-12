package store

import (
	"testing"
)

func TestPhoneParse(t *testing.T) {

	v := phoneParse("015030  001")
	if v != "5030" {
		t.Error("Expected 5030, got ", v)
	}

	v = phoneParse("010002 001")
	if v != "0002" {
		t.Error("Expected 0002, got ", v)
	}

	v = phoneParse("015088 001")
	if v != "5088" {
		t.Error("Expected 5088, got ", v)
	}

	v = phoneParse("015577 001")
	if v != "5577" {
		t.Error("Expected 5577, got ", v)
	}

	v = phoneParse("015079 001")
	if v != "5079" {
		t.Error("Expected 5079, got ", v)
	}

	v = phoneParse("015522 001")
	if v != "5522" {
		t.Error("Expected 5522, got ", v)
	}

	v = phoneParse("015522 002")
	if v != "015522 002" {
		t.Error("Expected 015522 002, got ", v)
	}

	v = phoneParse("015522002")
	if v != "015522002" {
		t.Error("Expected 015522002, got ", v)
	}

	v = phoneParse("9095564756")
	if v != "9095564756" {
		t.Error("Expected 9095564756, got ", v)
	}

	v = phoneParse("1155")
	if v != "1155" {
		t.Error("Expected 1155, got ", v)
	}

	v = phoneParse("")
	if v != "" {
		t.Error("Expected _, got ", v)
	}

	v = phoneParse(" ")
	if v != " " {
		t.Error("Expected space, got ", v)
	}
}
