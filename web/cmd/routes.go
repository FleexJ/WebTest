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

	router.HandleFunc("/logout/", app.logout)

	router.HandleFunc("/changeUser/", app.changeUserGET).Methods("GET")
	router.HandleFunc("/changeUser/", app.changeUserPOST).Methods("POST")

	router.HandleFunc("/changePassword/", app.changePasswordGET).Methods("GET")
	router.HandleFunc("/changePassword/", app.changePasswordPOST).Methods("POST")

	router.HandleFunc("/deleteUser/", app.deleteUserGET).Methods("GET")
	router.HandleFunc("/deleteUser/", app.deleteUserPOST).Methods("POST")

	return router
}
