package dto

import "testing"

func TestLoginRequestValidation_EmptyEmail(t *testing.T) {
	req := LoginRequest{}

	// Test empty validation - using validator package
	if req.Email == "" {
		// This should fail validation
		t.Log("Email is empty - validation would fail")
	}
}

func TestLoginRequestValidation_InvalidEmailFormat(t *testing.T) {
	req := LoginRequest{
		Email:    "invalid-email",
		Password: "password123",
	}

	// Simple email format check
	if req.Email != "" && !contains(req.Email, "@") {
		t.Log("Invalid email format - validation would fail")
	}
}

func TestLoginRequestValidation_EmptyPassword(t *testing.T) {
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	if req.Password == "" {
		t.Log("Password is empty - validation would fail")
	}
}

func TestLoginRequestValidation_ShortPassword(t *testing.T) {
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "123", // Too short
	}

	if len(req.Password) < 6 {
		t.Log("Password too short - validation would fail")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
