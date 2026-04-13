package utils

import "testing"

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{ErrInvalidInput, 400},
		{ErrUnauthorized, 401},
		{ErrForbidden, 403},
		{ErrNotFound, 404},
		{ErrInternal, 500},
	}

	for _, tt := range tests {
		code := GetStatusCode(tt.err)
		if code != tt.expected {
			t.Errorf("expected %d, got %d", tt.expected, code)
		}
	}
}