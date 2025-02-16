package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	type header struct {
		key   string
		value string
	}

	type returnValue struct {
		value string
		isErr bool
	}

	cases := []struct {
		input    header
		expected returnValue
	}{
		{
			input:    header{"Authorization", "Bearer validBearer"},
			expected: returnValue{"validBearer", false},
		},
		{
			input:    header{"Authorization", "Bearer "},
			expected: returnValue{"", true},
		},
		{
			input:    header{"Authorization", "BearerInvalidBearer"},
			expected: returnValue{"", true},
		},
		{
			input:    header{"Authorization", ""},
			expected: returnValue{"", true},
		},
		{
			input:    header{"Authorization", "Basic InvalidAuthType"},
			expected: returnValue{"", true},
		},
		{
			input:    header{"NotAValidHeader", "Bearer InvalidHeader"},
			expected: returnValue{"", true},
		},
		{
			input:    header{"", ""},
			expected: returnValue{"", true},
		},
	}

	for i, c := range cases {
		header := http.Header{}
		header.Set(c.input.key, c.input.value)
		result, err := GetBearerToken(header)
		if result != c.expected.value || (err != nil) != c.expected.isErr {
			t.Fatalf("\n[%v]\nExpected value: %v\nExpected error: %v\nResult value: %v\nResult has error:%v\n",
				i, c.expected.value, c.expected.isErr, result, err != nil)
		}
	}
}
