package main

import (
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

//Главная страница
func (app *application) indexPageGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.index.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		ts.Execute(w, struct {
			User *user
		}{
			User: nil,
		})
	} else {
		u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
		if err != nil {
			app.serverError(w, err)
			return
		}
		ts.Execute(w, struct {
			User *user
		}{
			User: u,
		})
	}
}

//Страница отображения всех пользователей
func (app *application) usersPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.users.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	users, err := getAllUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}

	u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, struct {
		User  *user
		Users []user
	}{
		User:  u,
		Users: users,
	})
}

//Отображение страницы регистрации
func (app *application) signUpPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.signUp.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, nil)
}

//Обработка POST-запроса страницы регистрации
func (app *application) signUpPagePOST(w http.ResponseWriter, r *http.Request) {
	u := user{
		Id:       bson.NewObjectId(),
		Email:    r.FormValue("email"),
		Name:     r.FormValue("name"),
		Surname:  r.FormValue("surname"),
		Password: r.FormValue("password"),
	}
	repPassword := r.FormValue("repPassword")

	if !u.valid(repPassword) {
		http.Redirect(w, r, "/signUp/", http.StatusSeeOther)
		return
	}
	err := u.saveUser()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("Новый пользователь:", u.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Отображение страницы авторизации
func (app *application) signInPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.signIn.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, nil)
}

//Обработка POST-запроса страницы авторизации
func (app *application) signInPagePOST(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Redirect(w, r, "/signIn/", http.StatusSeeOther)
		return
	}
	msg, err := auth(w, email, password)
	if msg != "" || err != nil {
		http.Redirect(w, r, "/signIn/", http.StatusSeeOther)
		return
	}
	app.infoLog.Println("Пользователь вошел:", email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Выход из учетной записи
func (app *application) logOut(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = tkn.deleteToken(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("Пользователь вышел:", tkn.EmailUser)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) changeUserGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.changeUser.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, struct {
		User   *user
		IdUser string
	}{
		User:   u,
		IdUser: u.Id.Hex(),
	})
}

//Обработка запроса на смену данных пользователя
func (app *application) changeUserPOST(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := user{
		Id:      bson.ObjectIdHex(tkn.IdUser),
		Email:   r.FormValue("email"),
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
	}

	uG, err := getUserById(u.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if uG == nil {
		http.Redirect(w, r, "/changeUser/", http.StatusSeeOther)
		return
	}
	u.Password = uG.Password
	if !u.valid(u.Password) {
		http.Redirect(w, r, "/changeUser/", http.StatusSeeOther)
		return
	}
	err = u.updateUser()
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Отображение страницы смены пароля
func (app *application) changePasswordGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.changePassword.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, struct {
		User   *user
		IdUser string
	}{
		User:   u,
		IdUser: u.Id.Hex(),
	})
}

//Обработка запроса на обновление пароля пользователя
func (app *application) changePasswordPOST(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(tkn.IdUser)
	password := r.FormValue("password")
	newPassword := r.FormValue("newPassword")
	repNewPassword := r.FormValue("repNewPassword")
	if password == "" || newPassword == "" || repNewPassword == "" ||
		newPassword != repNewPassword {
		http.Redirect(w, r, "/changePassword/", http.StatusSeeOther)
		return
	}
	u, err := getUserById(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u == nil {
		http.Redirect(w, r, "/changePassword/", http.StatusSeeOther)
		return
	}
	err = u.comparePassword(password)
	if err != nil {
		http.Redirect(w, r, "/changePassword/", http.StatusSeeOther)
		return
	}
	u.Password = newPassword
	err = u.updateUserPassword()
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteUserGET(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := getUserById(bson.ObjectIdHex(tkn.IdUser))
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	ts, err := template.ParseFiles(
		"./ui/views/page.deleteUser.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts.Execute(w, struct {
		User   *user
		IdUser string
	}{
		User:   u,
		IdUser: u.Id.Hex(),
	})
}

func (app *application) deleteUserPOST(w http.ResponseWriter, r *http.Request) {
	tkn, err := checkAuth(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(tkn.IdUser)
	email := r.FormValue("email")
	password := r.FormValue("password")
	if password == "" || email == "" {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}
	u, err := getUserById(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u == nil {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}
	err = u.comparePassword(password)
	if err != nil {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}
	if email != u.Email {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}
	err = u.deleteUser()
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = tkn.deleteToken(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/logout/", http.StatusSeeOther)
}
