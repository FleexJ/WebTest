package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", app.indexPageGET).Methods("GET")

	router.HandleFunc("/users/", app.usersPageGET).Methods("GET")

	router.HandleFunc("/signUp/", app.signUpPageGET).Methods("GET")
	router.HandleFunc("/signUp/", app.signUpPagePOST).Methods("POST")

	router.HandleFunc("/signIn/", app.signInPageGET).Methods("GET")
	router.HandleFunc("/signIn/", app.signInPagePOST).Methods("POST")

	router.HandleFunc("/logout/", app.logOut)

	router.HandleFunc("/changeUser/", app.changeUserGET).Methods("GET")
	router.HandleFunc("/changeUser/", app.changeUserPOST).Methods("POST")

	//TODO добавить изменение пароля

	return router
}
