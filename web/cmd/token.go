package main

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

type token struct {
	IdUser    string
	EmailUser string
	Token     string
	Expires   int64
}

//Проверка токена на пустоту
func (t token) isEmpty() bool {
	if t.EmailUser == "" || t.Token == "" || t.IdUser == "" || t.Expires == 0 {
		return true
	}
	return false
}

//Сохраняет токен в базе
func (t token) saveToken(w http.ResponseWriter) error {
	http.SetCookie(w,
		newCookie(idCookieName, t.IdUser))
	http.SetCookie(w,
		newCookie(emailCookieName, t.EmailUser))
	http.SetCookie(w,
		newCookie(tokenCookieName, t.Token))
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
		newCookie(emailCookieName, ""))
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
	_, err = collection.RemoveAll(bson.M{"id": t.IdUser, "token": t.Token})
	if err != nil {
		return err
	}
	return nil
}

//Проверка на существование токена в базе
func (t token) findInDB() (bool, error) {
	session, err := getSession()
	if err != nil {
		return false, err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	var tkn *token
	err = collection.Find(bson.M{"iduser": t.IdUser, "token": t.Token}).One(&tkn)
	if err != nil && err.Error() == "not found" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
