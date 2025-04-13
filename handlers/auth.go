package handlers

import (
	"go-web-app/db"
	"go-web-app/models"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	Store = s
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	switch r.Method {
	case "GET":
		tmpl.Execute(w, nil)
	case "POST":
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.User
		result := db.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			tmpl.Execute(w, map[string]string{"Error": "User Not exists"})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			tmpl.Execute(w, map[string]string{
				"Error": "Invalid credentials",
			})
			return
		}
		session, _ := Store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
