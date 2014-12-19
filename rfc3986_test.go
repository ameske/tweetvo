package main

import "testing"

var percentEncodeTests = []struct {
	in  string
	out string
}{
	{"Ladies + Gentlemen", "Ladies%20%2B%20Gentlemen"},
	{"An encoded string!", "An%20encoded%20string%21"},
	{"Dogs, Cats & Mice", "Dogs%2C%20Cats%20%26%20Mice"},
	{"☃", "%E2%98%83"},
}

func TestPercentEncode(t *testing.T) {
	for _, test := range percentEncodeTests {
		if got := percentEncode(test.in); got != test.out {
			t.Errorf("Got: %s\t Want: %s", got, test.out)
		}
	}
}
