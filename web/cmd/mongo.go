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

//Получение пользователя по адресу почты
func getUserByEmail(email string) *User {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	var u User
	err = collection.Find(bson.M{"email": email}).One(&u)
	if err != nil {
		return nil
	}

	return &u
}

func getUserById(id bson.ObjectId) *User {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	var u User
	err = collection.Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		return nil
	}

	return &u
}

//Возвращает список всех пользователей
func getAllUsers() []User {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	var users []User
	err = collection.Find(bson.M{}).All(&users)
	if err != nil {
		return nil
	}

	return users
}

//Удаляет все токены из базы
func deleteAllOldTokens() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(authCol)
	_, err = collection.RemoveAll(nil)
	if err != nil {
		return err
	}

	return nil
}
