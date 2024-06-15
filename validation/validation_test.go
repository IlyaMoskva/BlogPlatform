package validation

import (
	"testing"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  int
		shouldErr bool
	}{
		{
			name:      "Valid ID",
			input:     "1",
			expected:  1,
			shouldErr: false,
		},
		{
			name:      "Valid ID with multiple digits",
			input:     "123",
			expected:  123,
			shouldErr: false,
		},
		{
			name:      "Invalid ID with negative number",
			input:     "-1",
			expected:  0,
			shouldErr: true,
		},
		{
			name:      "Invalid ID with zero",
			input:     "0",
			expected:  0,
			shouldErr: true,
		},
		{
			name:      "Invalid ID with non-numeric characters",
			input:     "abc",
			expected:  0,
			shouldErr: true,
		},
		{
			name:      "Empty ID",
			input:     "",
			expected:  0,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateID(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateID() error = %v, shouldErr %v", err, tt.shouldErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ValidateID() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
