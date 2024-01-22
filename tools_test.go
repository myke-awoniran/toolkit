package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("expect length to be 10, got", len(s))
	}
}
