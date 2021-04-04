package main

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

type token struct {
	IdUser  string
	Token   string
	Expires int64
}

//Проверка токена на пустоту
func (t token) isEmpty() bool {
	if t.Token == "" || t.IdUser == "" || t.Expires == 0 {
		return true
	}
	return false
}

//Сохраняет токен в базе
func (t token) saveToken(w http.ResponseWriter) error {
	http.SetCookie(w,
		newCookie(idCookieName, t.IdUser))
	//base64 token save in cookie
	base64Tkn := base64.StdEncoding.EncodeToString([]byte(t.Token))
	http.SetCookie(w,
		newCookie(tokenCookieName, base64Tkn))
	http.SetCookie(w,
		newCookie(expiresCookieName, strconv.FormatInt(t.Expires, 10)))
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	//Удаление всех устаревших токенов
	var tkns []token
	_ = collection.Find(bson.M{"id": t.IdUser}).All(&tkns)
	for _, el := range tkns {
		if el.Expires <= time.Now().Unix() {
			_ = collection.Remove(bson.M{"expires": el.Expires, "iduser:": el.IdUser})
		}
	}
	//bcrypt token save in DB
	bcryptTkn, err := bcrypt.GenerateFromPassword([]byte(t.Token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.Token = string(bcryptTkn)
	err = collection.Insert(t)
	if err != nil {
		return err
	}
	return nil
}

//Удаляет токен из базы и из куки
func (t token) deleteToken(w http.ResponseWriter) error {
	http.SetCookie(w,
		newCookie(idCookieName, ""))
	http.SetCookie(w,
		newCookie(tokenCookieName, ""))
	http.SetCookie(w,
		newCookie(expiresCookieName, ""))
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	var tkns []token
	//Считываем из базы все токены текущего пользователя
	err = collection.Find(bson.M{"iduser": t.IdUser}).All(&tkns)
	if err != nil && err.Error() != "not found" {
		return err
	}
	//Декодируем токен из куки
	tDecode, err := base64.StdEncoding.DecodeString(t.Token)
	if err != nil {
		return err
	}
	for _, tkn := range tkns {
		if bcrypt.CompareHashAndPassword([]byte(tkn.Token), tDecode) == nil {
			//Удаление токена из БД
			_, err = collection.RemoveAll(bson.M{"iduser": tkn.IdUser, "token": tkn.Token})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//Проверка на существование токена в базе
func (t token) findInDB() bool {
	session, err := getSession()
	if err != nil {
		return false
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	var tkns []token
	//Считываем из базы все токены текущего пользователя
	err = collection.Find(bson.M{"iduser": t.IdUser}).All(&tkns)
	if err != nil && err.Error() == "not found" {
		return false
	}
	if err != nil {
		return false
	}
	//Декодируем токен из куки
	tDecode, err := base64.StdEncoding.DecodeString(t.Token)
	if err != nil {
		return false
	}
	//Ищем совпадения токена в куки и БД
	for _, tkn := range tkns {
		if bcrypt.CompareHashAndPassword([]byte(tkn.Token), tDecode) == nil {
			return true
		}
	}
	return false
}
