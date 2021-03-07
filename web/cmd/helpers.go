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

//Вощвращает токен, считанный из куки
func getTokenCookies(r *http.Request) *token {
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

//Возвращает новый объект куки
func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
}

//Функция очистки куки
func clearCookies(w http.ResponseWriter) {
	http.SetCookie(w,
		newCookie(emailCookieName, ""))
	http.SetCookie(w,
		newCookie(tokenCookieName, ""))
}

//Функция авторизации пользователя
//Ищет совпадения в базе пользователей
//Выдает новый токен доступа
//при успехе возвращается пустая строка
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
	http.SetCookie(w,
		newCookie(emailCookieName, token.EmailUser))
	http.SetCookie(w,
		newCookie(tokenCookieName, token.Token))
	return "", nil
}

//Генерирует новый токен на основе почты
func generateToken(word string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return word + string(b)
}

//Проверка токена доступа, возвращает токен с данными при успехе
func checkAuth(r *http.Request) *token {
	token := getTokenCookies(r)
	if token.isEmpty() {
		return nil
	}
	is := token.findInDB()
	if !is {
		return nil
	}
	return token
}
