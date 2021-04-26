package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

const regexEmail = `^\w+@\w+[.]\w+$`

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Email    string
	Name     string
	Surname  string
	Password string
}

//Валидация пользователя перед записью в базу
func (u *User) valid(repPassword string) bool {
	matched, _ := regexp.MatchString(regexEmail, u.Email)
	if !matched ||
		u.Name == "" ||
		u.Surname == "" ||
		u.Password == "" ||
		u.Password != repPassword {
		return false
	}

	uG := getUserByEmail(u.Email)
	if uG != nil && u.Id != uG.Id {
		return false
	}

	return true
}

//Сравнение пароля пользователя
func (u *User) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}

//Сохранение пользователя в базе
func (u User) saveUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(bcryptPassw)
	u.Id = bson.NewObjectId()
	err = collection.Insert(u)
	if err != nil {
		return err
	}

	return nil
}

//Обновление данных пользователя
func (u User) updateUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Update(bson.M{"_id": u.Id}, u)
	if err != nil {
		return err
	}

	return nil
}

//Обновление пароля пользователя
func (u User) updateUserPassword(password string) error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(bcryptPassw)
	err = collection.Update(bson.M{"_id": u.Id}, u)
	if err != nil {
		return err
	}

	return nil
}

//Удаление пользователя из базы
func (u User) deleteUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Remove(bson.M{"_id": u.Id})
	if err != nil {
		return err
	}

	return nil
}
