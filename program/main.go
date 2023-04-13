package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var tpl *template.Template
var admin = user{
	username: "anas",
	password: "8590"}
var currentUser string 


// create a new session store with a secret key
var store = sessions.NewCookieStore([]byte("secret-key"))


func main() {

	tpl, _ = template.ParseGlob("*.html")
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/login", submitHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/signOut", signOutHandler)
	http.ListenAndServe(":9999", nil)

}

//	http.HandleFunc("/", loginHandler)
//
// it will take the session and evalute if the user exist then redirect to home page else return to loginpage
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", " no-store,no-cache,must-revalidate")

	session, err := store.Get(r, "session-name")
	errorHandling(err)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {

		http.Redirect(w, r, "/home", http.StatusFound)
		

		return
	}
	w.Header().Set("Cache-Control", " no-store,no-cache,must-validate")

	err = tpl.ExecuteTemplate(w, "index.html", nil)
	errorHandling(err)
}

//	http.HandleFunc("/login", submitHandler)
//
// index handler will perform when the submit button press and create a session for the permitted user
func submitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", " no-store,no-cache,must-revalidate")


	session, err := store.Get(r, "session-name")
	errorHandling(err)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	r.ParseForm()

	currentUser := r.FormValue("username")
	password := r.FormValue("password")

	if currentUser != admin.username || password != admin.password {


		tpl.ExecuteTemplate(w, "index.html", "Invalid username and password")
		return

	}
	session.Values["authenticated"] = true

	err = session.Save(r, w)
	errorHandling(err)
	err = tpl.ExecuteTemplate(w, "home.html", currentUser)
	errorHandling(err)
}

//	http.HandleFunc("/home", homeHandler)
//
// home handler evalute the user if the user exist  if the user !exist then redirect to the login page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", " no-store,no-cache,must-revalidate")


	session, err := store.Get(r, "session-name")
	errorHandling(err)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = tpl.ExecuteTemplate(w, "home.html", admin.username)
	errorHandling(err)
}

//	http.HandleFunc("/signOut", signOutHandler)=
//
// This handler will delete the session and logout from the site
func signOutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	errorHandling(err)
	session.Values["authenticated"] = false
	err = session.Save(r, w)
	errorHandling(err)
	
	w.Header().Set("Cache-Control", "no-cache , no-store , must-revalidate")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	

}

type user struct {
	username string
	password string
}

func errorHandling(err error) {

	if err != nil {
		fmt.Println("there is an error", err)

	}
}
