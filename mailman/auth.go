package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const authCookie = "gms-auth-token"
const loggedInCookie = "gms-auth-loggedIn"

// LoginRequest is used to unmarshal the simple login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(ctx context.Context, headers map[string]string, body string) (map[string][]string, error) {
	loginCreds := LoginRequest{}

	err := json.Unmarshal([]byte(body), &loginCreds)
	if err != nil {
		return nil, err
	}

	// TODO: Verify credentials
	if loginCreds.Username != "gideonw" || loginCreds.Password != "test" {
		return nil, fmt.Errorf("Error: passwrod incorrect for %s", loginCreds.Username)
	}

	ret := make(map[string][]string)
	setCookieHeaders := []string{}

	token, err := signNewToken(loginCreds.Username)
	if err != nil {
		return nil, err
	}

	authToken := http.Cookie{
		Name:   authCookie,
		Domain: "gps." + domain,
		Path:   "/",
		Value:  token,

		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Secure:   true,
		MaxAge:   1209600, // 14 days
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	setCookieHeaders = append(setCookieHeaders, authToken.String())

	loggedIn := http.Cookie{
		Name:   loggedInCookie,
		Domain: "gps." + domain,
		Path:   "/",
		Value:  "true",

		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Secure:   true,
		MaxAge:   1209600, // 14 days
		SameSite: http.SameSiteStrictMode,
	}
	setCookieHeaders = append(setCookieHeaders, loggedIn.String())

	ret["Set-Cookie"] = setCookieHeaders

	return ret, nil
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