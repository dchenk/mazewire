package base58

import (
	"testing"
)

func TestAlphaLen(t *testing.T) {
	if len(alphabet) != 58 {
		t.Errorf("the sky is falling")
	}
}
