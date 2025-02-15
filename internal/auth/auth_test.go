package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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

func TestJWT(t *testing.T) {
	cases := []struct {
		id        uuid.UUID
		secret    string
		expiresAt time.Duration
	}{
		{
			id:        uuid.New(),
			secret:    "Test",
			expiresAt: time.Hour,
		},
		{
			id:        uuid.New(),
			secret:    "ttt",
			expiresAt: time.Minute,
		},
	}

	for i, c := range cases {
		t.Logf("Case %v: %v, %v, %v\n", i, c.id, c.secret, c.expiresAt)
		token, err := MakeJWT(c.id, c.secret, c.expiresAt)
		if err != nil {
			t.Fatalf("Error creating jwt token: %v, %v, %v - %v", c.id, c.secret, c.expiresAt, err)
		}

		t.Logf("Token: %v\n", token)
		id, err := ValidateJWT(token, c.secret)
		if err != nil {
			t.Fatalf("Error validating token: %v, %v - %v", token, c.secret, err)
		}

		t.Logf("ID: %v\n", id)
		if c.id != id {
			t.Fatalf("Id's doesn't match: %v != %v", c.id, id)
		}
	}
}
