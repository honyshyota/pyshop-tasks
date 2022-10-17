package checkEven

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckEven_isEven test function even from package checkEven
func TestCheckEven_isEven(t *testing.T) {
	// Iterate over even numbers
	for i := 0; i < 100; i += 2 { // from 0 to 100
		checkEven := isEven(i)
		// use package testify, function Equal() сhecks if two variables match
		assert.Equal(t, false, checkEven)
	}
	// Iterate over odd numbers
	for i := 1; i < 100; i += 2 {
		checkEven := isEven(i)
		assert.Equal(t, true, checkEven)
	}
}
