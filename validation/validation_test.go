package validation

import (
	"net/http"
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

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectError  bool
		expectedMsg  string
		expectedCode int
	}{
		{
			name:        "Valid query",
			query:       "test",
			expectError: false,
		},
		{
			name:        "Valid query with leading and trailing spaces",
			query:       "  test  ",
			expectError: false,
		},
		{
			name:         "Empty query",
			query:        "",
			expectError:  true,
			expectedMsg:  "query parameter is required",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Query too long",
			query:        string(make([]byte, 101)), // 101 characters long
			expectError:  true,
			expectedMsg:  "query parameter is too long",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid characters",
			query:        "!@#$%^&*()",
			expectError:  true,
			expectedMsg:  "query parameter contains invalid characters",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Non-UTF-8 characters",
			query:        string([]byte{0xff, 0xfe, 0xfd}),
			expectError:  true,
			expectedMsg:  "query parameter is not valid UTF-8",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.query)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateQuery() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if err != nil {
				if httpErr, ok := err.(HttpError); ok {
					if httpErr.Message != tt.expectedMsg {
						t.Errorf("ValidateQuery() error message = %v, expectedMsg %v", httpErr.Message, tt.expectedMsg)
					}
					if httpErr.Code != tt.expectedCode {
						t.Errorf("ValidateQuery() error code = %v, expectedCode %v", httpErr.Code, tt.expectedCode)
					}
				} else {
					t.Errorf("ValidateQuery() error is not of type HttpError")
				}
			}
		})
	}
}
