package safeconverter

import "testing"

func TestToInt(t *testing.T) {
	if ToInt(uint64(12)) != 12 {
		t.Fatalf("expected 12")
	}
}

func TestToInt64(t *testing.T) {
	if ToInt64(int32(34)) != 34 {
		t.Fatalf("expected 34")
	}
}
