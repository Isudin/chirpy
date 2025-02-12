package auth

import "testing"

func TestHashPassword(t *testing.T) {
	passwords := []string{"test123", "ADF123!@#fs"}
	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatal("Error hashing password")
		}

		if CheckPasswordHash(password, hash) != nil {
			t.Fatalf("Password '%v' doesn't match the hash '%v'", password, hash)
		}
	}
}
