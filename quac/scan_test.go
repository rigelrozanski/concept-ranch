package quac

import "testing"

func TestEquals(t *testing.T) {
	c1 := NewColour(1000, 1000, 1000)
	c2 := NewColour(1000, 1000, 1000)
	if !c1.Equals(c2) {
		t.Errorf("should have been equal\n\tc1: %s\n\tc2: %s", c1, c2)
	}
}
