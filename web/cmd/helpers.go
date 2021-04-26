package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	idCookieName    = "id"
	tokenCookieName = "token"
	expDay          = 60 * 24
)

//Возвращает токен, считанный из куки
func getTokenCookies(r *http.Request) *Token {
	cookieId, err := r.Cookie(idCookieName)
	if err != nil {
		return nil
	}

	cookieToken, err := r.Cookie(tokenCookieName)
	if err != nil {
		return nil
	}

	if cookieId.Value == "" || cookieToken.Value == "" {
		return nil
	}

	return &Token{
		IdUser: cookieId.Value,
		Token:  cookieToken.Value,
	}
}

//Возвращает новый объект куки
func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(expDay * time.Hour),
	}
}

//Функция авторизации пользователя
//Ищет совпадения в базе пользователей
//Выдает новый токен доступа
//при успехе нет ошибки
func auth(w http.ResponseWriter, email, password string) error {
	u := getUserByEmail(email)
	if u == nil {
		return errors.New("user not found")
	}

	err := u.comparePassword(password)
	if err != nil {
		return err
	}

	genToken, err := generateToken(u.Id.Hex())
	if err != nil {
		return err
	}

	tkn := Token{
		IdUser: u.Id.Hex(),
		Token:  genToken,
	}

	err = tkn.saveToken(w)
	if err != nil {
		return err
	}
	return nil
}

//Генерирует новый токен на основе слова
func generateToken(word string) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	bcryptB, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return word + strconv.FormatInt(time.Now().Unix(), 10) + string(bcryptB), nil
}

//Проверка токена доступа, возвращает токен с данными и текущего пользователя при успехе
func checkAuth(r *http.Request) (*Token, *User) {
	tkn := getTokenCookies(r)
	if tkn == nil {
		return nil, nil
	}

	is := tkn.findInDB()
	if !is {
		return nil, nil
	}

	u := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if u == nil {
		return nil, nil
	}
	return tkn, u
}
