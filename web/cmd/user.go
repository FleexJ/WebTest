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
func (usr *User) valid(repPassword string) bool {
	matched, _ := regexp.MatchString(regexEmail, usr.Email)
	if !matched ||
		usr.Name == "" ||
		usr.Surname == "" ||
		usr.Password == "" ||
		usr.Password != repPassword {
		return false
	}

	uG := getUserByEmail(usr.Email)
	if uG != nil && usr.Id != uG.Id {
		return false
	}

	return true
}

//Сравнение пароля пользователя
func (usr *User) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}

//Сохранение пользователя в базе
func (usr User) saveUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(bcryptPassw)
	usr.Id = bson.NewObjectId()
	err = collection.Insert(usr)
	if err != nil {
		return err
	}

	return nil
}

//Обновление данных пользователя
func (usr User) updateUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Update(bson.M{"_id": usr.Id}, usr)
	if err != nil {
		return err
	}

	return nil
}

//Обновление пароля пользователя
func (usr User) updateUserPassword(password string) error {
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

	usr.Password = string(bcryptPassw)
	err = collection.Update(bson.M{"_id": usr.Id}, usr)
	if err != nil {
		return err
	}

	return nil
}

//Удаление пользователя из базы
func (usr User) deleteUser() error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Remove(bson.M{"_id": usr.Id})
	if err != nil {
		return err
	}

	return nil
}
