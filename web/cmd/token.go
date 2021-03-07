package main

import (
	"gopkg.in/mgo.v2/bson"
)

type token struct {
	EmailUser string
	Token     string
}

//Проверка токена на пустоту
func (t token) isEmpty() bool {
	if t.EmailUser == "" || t.Token == "" {
		return true
	}
	return false
}

//Сохраняет токен в базе
func (t token) saveToken() error {
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

//Удаляет токен из базы
func (t token) deleteToken() error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	_, err = collection.RemoveAll(bson.M{"emailuser": t.EmailUser, "token": t.Token})
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
	var token token
	err = collection.Find(bson.M{"emailuser": t.EmailUser, "token": t.Token}).One(&token)
	if err != nil {
		return false
	}
	if token.isEmpty() {
		return false
	}
	return true
}
