package main

import (
	"gopkg.in/mgo.v2/bson"
)

type token struct {
	EmailUser string
	Token     string
}

func (t token) isEmpty() bool {
	if t.EmailUser == "" || t.Token == "" {
		return true
	}
	return false
}

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

func (t token) findInDB() (bool, error) {
	session, err := getSession()
	if err != nil {
		return false, err
	}
	defer session.Close()
	collection := session.DB(database).C(authCol)
	var token token
	err = collection.Find(bson.M{"emailuser": t.EmailUser, "token": t.Token}).One(&token)
	if err != nil {
		return false, err
	}
	if token.isEmpty() {
		return false, nil
	}
	return true, nil
}
