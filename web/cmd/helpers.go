package main

import (
	"encoding/base64"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/http"
	"time"
)

const (
	idCookieName    = "id"
	emailCookieName = "email"
	tokenCookieName = "token"
)

//Вощвращает токен, считанный из куки
func getTokenCookies(r *http.Request) *token {
	cookieId, err := r.Cookie(idCookieName)
	if err != nil {
		return nil
	}
	cookieEmail, err := r.Cookie(emailCookieName)
	if err != nil {
		return nil
	}
	cookieToken, err := r.Cookie(tokenCookieName)
	if err != nil {
		return nil
	}
	if cookieId.Value == "" || cookieToken.Value == "" || cookieEmail.Value == "" {
		return nil
	}
	return &token{
		IdUser:    cookieId.Value,
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
	tkn := token{
		IdUser:    u.Id.Hex(),
		EmailUser: u.Email,
		Token:     generateToken(u.Id.Hex()),
	}
	err = tkn.saveToken(w)
	if err != nil {
		return "", err
	}
	return "", nil
}

//Генерирует новый токен на основе какого-то слова
func generateToken(word string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return word + base64.StdEncoding.EncodeToString(b)
}

//Проверка токена доступа, возвращает токен с данными при успехе
func checkAuth(r *http.Request) (*token, error) {
	tkn := getTokenCookies(r)
	if tkn == nil {
		return nil, nil
	}
	is, err := tkn.findInDB()
	if err != nil {
		return nil, err
	}
	if !is {
		return nil, nil
	}
	u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return tkn, nil
}
