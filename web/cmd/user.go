package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

const regexEmail = `^\w+@\w+[.]\w+$`

type user struct {
	Id       bson.ObjectId `bson:"_id"`
	Email    string
	Name     string
	Surname  string
	Password string
}

//Валидация пользователя перед записью в базу
func (u *user) valid(repPassword string) bool {
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

//Проверка пользователя на пустоту
func (u user) isEmpty() bool {
	empty := user{}
	if u == empty {
		return true
	}
	return false
}

//Сравнение пароля пользователя
func (u *user) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return nil
	}
	return err
}

//Сохранение пользователя в базе
func (u user) saveUser() error {
	session, err := getSession()
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
func (u user) updateUser() error {
	session, err := getSession()
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
func (u user) updateUserPassword(password string) error {
	session, err := getSession()
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
func (u user) deleteUser() error {
	session, err := getSession()
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
