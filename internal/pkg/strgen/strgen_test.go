package strgen

import (
	"testing"

	"github.com/GeorgeGogos/goaap/pkg/util"
)

// TestRandHexString verifies the accuracy of the hexadecimal string generator.
func TestRandHexString(t *testing.T) {
	length := 16
	randStr := util.RandHexString(16)

	if len(randStr) != length {
		t.Errorf("Expected length %d, got %d", length, len(randStr))
	}

	for _, char := range randStr {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("Invalid character in hex string: %c", char)
		}
	}
}

// BenchmarkRandHexString measures the resource usage of the hexadecimal string generator.
func BenchmarkRandHexString(b *testing.B) {
	length := 16
	for i := 0; i < b.N; i++ {
		util.RandHexString(length)
	}
}
