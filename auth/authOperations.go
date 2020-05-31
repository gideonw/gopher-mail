package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

const authCookie = "gms-auth-token"
const loggedInCookie = "gms-auth-loggedIn"

// LoginRequest is used to unmarshal the simple login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login ...
func Login(ctx context.Context, domain string, headers map[string]string, body string) (TokenPayload, error) {
	loginCreds := LoginRequest{}

	err := json.Unmarshal([]byte(body), &loginCreds)
	if err != nil {
		return TokenPayload{}, err
	}

	// TODO: Verify credentials
	if loginCreds.Username != "gideonw" || loginCreds.Password != "test" {
		return TokenPayload{}, fmt.Errorf("Error: passwrod incorrect for %s", loginCreds.Username)
	}

	token, err := signNewToken(loginCreds.Username)
	if err != nil {
		return TokenPayload{}, err
	}

	return TokenPayload{
		Token: token,
	}, nil
}

func signNewToken(username string) (string, error) {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("%v", ss)

	return ss, nil
}

// WellKnownOpenIDConfig OpenID standard config
func WellKnownOpenIDConfig() (string, error) {
	iss := "https://gps.gideonw.xyz"
	openID := OpenIDConfig{
		Issuer:                 iss,
		AuthorizationEndpoint:  iss + "/api/auth/authorize",
		TokenEndpoint:          iss + "/api/auth/token",
		JWKSURI:                iss + "/api/auth/jwks.json",
		ResponseTypesSupported: []string{},
	}

	buf, err := json.Marshal(openID)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// WellKnownJWKSJSON OpenID standard public JWT signature keys
func WellKnownJWKSJSON() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	pubJWK, err := jwk.New(key.Public())
	if err != nil {
		return "", err
	}

	jwks := JWKS{
		Keys: []jwk.Key{
			pubJWK,
		},
	}

	buf, err := json.Marshal(jwks)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
