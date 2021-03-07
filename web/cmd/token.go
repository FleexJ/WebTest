package main

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type token struct {
	IdUser    string
	EmailUser string
	Token     string
}

//Проверка токена на пустоту
func (t token) isEmpty() bool {
	if t.EmailUser == "" || t.Token == "" || t.IdUser == "" {
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
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
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
func (t token) findInDB() bool {
	session, err := getSession()
	if err != nil {
		return false
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	var tkn *token
	err = collection.Find(bson.M{"iduser": t.IdUser, "token": t.Token}).One(&tkn)
	if err != nil {
		return false
	}
	if tkn == nil {
		return false
	}
	return true
}
