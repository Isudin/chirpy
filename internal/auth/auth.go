package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issuedAt := jwt.NumericDate{
		Time: time.Now().UTC(),
	}

	expiresAt := jwt.NumericDate{
		Time: time.Now().UTC().Add(expiresIn),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  &issuedAt,
			ExpiresAt: &expiresAt,
			Subject:   userID.String(),
		},
		nil,
	)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	idStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	return getToken(headers, "Bearer")
}

func GetAPIKey(headers http.Header) (string, error) {
	return getToken(headers, "ApiKey")
}

func getToken(headers http.Header, tokenType string) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("no authorization header found")
	}

	tokenString, isFound := strings.CutPrefix(token, tokenType+" ")
	if !isFound {
		return "", errors.New("invalid authorization header type")
	}

	if tokenString == "" {
		return "", errors.New("no bearer token found")
	}

	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 256)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
