package http

import (
	"reflect"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	timestamp := "1679554499156"
	raw := "Abcd1234"

	out, err := EncryptPassword(raw, timestamp)

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "EoknC2O10yGB7wVowXxiCPdbvH+aU35Nvted5lmEKzXPc2wHRzXS/6yxJoRz4igc", out)
}

func assertEqual(t *testing.T, a, b interface{}) {
	ok := reflect.DeepEqual(a, b)

	if !ok {
		t.Errorf("not equal, expect: %v, actual: %v", a, b)
		t.FailNow()
		return
	}
}
