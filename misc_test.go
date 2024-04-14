package bilibili

import "testing"

func TestAvBv(t *testing.T) {
	if Av2Bv(111298867365120) != "BV1L9Uoa9EUx" {
		t.Fail()
	}
	if 111298867365120 != Bv2Av("BV1L9Uoa9EUx") {
		t.Fail()
	}
}
