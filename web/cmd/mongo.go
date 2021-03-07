package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoUrl = "mongodb://localhost:27017"
	database = "web"
	usersCol = "users"
	authCol  = "aut"
)

//Получение сессии для связи с базой
func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	return session, nil
}

//Получение пользователя под адресу почты
func getUserByEmail(email string) (*user, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(database).C(usersCol)
	var u user
	err = collection.Find(bson.M{"email": email}).One(&u)
	if err != nil && err.Error() == "not found" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//Возвращается список всех пользователей
func getAllUsers() ([]user, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(database).C(usersCol)
	var users []user
	err = collection.Find(bson.M{}).All(&users)
	if err != nil && err.Error() == "not found" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}
