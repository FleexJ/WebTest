package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	id bson.ObjectId `bson:"_id"`
	Email string
	Name string
	Surname string
	Password string
}

func (u *user) valid(repPassword string) bool {
	if u.Email == "" ||
		u.Name == "" ||
		u.Surname == "" ||
		u.Password == "" ||
		u.Password != repPassword {
		return false
	}

	if getUserByEmail(u.Email) != nil{
		return false
	}

	return true
}

func (u user) isEmpty() bool {
	empty := user{}
	if u == empty {
		return true
	}
	return false
}

func (u *user) comparePassword(password string)  error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}


func (u user) saveUser() error {
	session, err := getSession()
	defer session.Close()
	if err != nil {
		return err
	}
	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bcryptPassw)
	err = collection.Insert(u)
	if err != nil {
		return err
	}

	return nil
}
