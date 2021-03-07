package main

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	emailCookieName = "email"
	tokenCookieName = "token"
)

func getCookies(r *http.Request) *token {
	cookieEmail, err := r.Cookie(emailCookieName)
	if err != nil {
		return &token{}
	}
	cookieToken, err := r.Cookie(tokenCookieName)
	if err != nil {
		return &token{}
	}
	return &token{
		EmailUser: cookieEmail.Value,
		Token:     cookieToken.Value,
	}
}

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
}

func clearCookies(w http.ResponseWriter) {
	cookieId := newCookie(emailCookieName, "")
	cookieToken := newCookie(tokenCookieName, "")
	http.SetCookie(w, cookieToken)
	http.SetCookie(w, cookieId)
}

func auth(w http.ResponseWriter, email, password string) (string, error) {
	u, err := getUserByEmail(email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "User not found", nil
	}

	err = u.comparePassword(password)
	if err != nil {
		return "", err
	}

	token := token{
		EmailUser: u.Email,
		Token:     generateToken(u.Email),
	}
	err = token.saveToken()
	if err != nil {
		return "", err
	}

	cookieId := newCookie(emailCookieName, token.EmailUser)
	cookieToken := newCookie(tokenCookieName, token.Token)
	http.SetCookie(w, cookieId)
	http.SetCookie(w, cookieToken)

	return "", nil
}

func generateToken(word string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return word + string(b)
}

func checkAuth(r *http.Request) (*token, error) {
	token := getCookies(r)
	if token.isEmpty() {
		return nil, nil
	}
	is, err := token.findInDB()
	if err != nil {
		return nil, err
	}
	if !is {
		return nil, nil
	}
	return token, nil
}
