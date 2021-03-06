package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	mux := mux.NewRouter()

	mux.HandleFunc("/", app.indexPageGET).Methods("GET")

	mux.HandleFunc("/users/", app.usersPageGET).Methods("GET")

	mux.HandleFunc("/signUp/", app.signUpPageGET).Methods("GET")
	mux.HandleFunc("/signUp/", app.signUpPagePOST).Methods("POST")

	mux.HandleFunc("/signIn/", app.signInPageGET).Methods("GET")
	mux.HandleFunc("/signIn/", app.signInPagePOST).Methods("POST")

	mux.HandleFunc("/logout/", app.logOut)

	return mux
}
