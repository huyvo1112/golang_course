package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeGenerator(t *testing.T) {
	t.Parallel()

	testcase := []struct {
		name   string
		length int
	}{
		{
			name:   "Normal case: Test code generation with length 8",
			length: 8,
		},
		{
			name:   "Normal case:Test code generation with length 12",
			length: 12,
		},
		{
			name:   "Edge case: Test code generation with valid length 0",
			length: 0,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewCodeGenerator()
			code, err := generator.GenerateCode(tc.length)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			assert.Equal(t, tc.length, len(code), "Generated code length does not match expected length")
		})
	}
}
