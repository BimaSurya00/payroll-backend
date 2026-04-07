package helper

import "testing"

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"
	hash, _ := HashPassword(password)

	// Test correct password
	if !CheckPassword(hash, password) {
		t.Error("Should match correct password")
	}

	// Test wrong password
	wrongPassword := "WrongPassword"
	if CheckPassword(hash, wrongPassword) {
		t.Error("Should not match wrong password")
	}
}

func TestToday(t *testing.T) {
	date := Today()

	if date.IsZero() {
		t.Error("Today should return a valid date")
	}
}

func TestNow(t *testing.T) {
	now := Now()

	if now.IsZero() {
		t.Error("Now should return a valid time")
	}
}
