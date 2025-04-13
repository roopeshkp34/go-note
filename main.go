package main

import (
	"fmt"
	"go-web-app/db"
	"go-web-app/middleware"
	"go-web-app/models"

	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

var store = middleware.Store

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("templates", tmpl)
	t := template.Must(template.ParseFiles(tmplPath))
	t.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, "login.html", nil)
	case "POST":
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user models.User
		result := db.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			renderTemplate(w, "login.html", map[string]string{"Error": "User Not exists"})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			renderTemplate(w, "login.html", map[string]string{
				"Error": "Invalid credentials",
			})
			return
		}
		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	db.DB.Find(&notes)
	renderTemplate(w, "dashboard.html", notes)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	note := models.Note{Title: title, Content: content}
	db.DB.Create(&note)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func editNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var note models.Note
	db.DB.First(&note, id)
	note.Title = title
	note.Content = content
	db.DB.Save(&note)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	db.DB.Delete(&models.Note{}, id)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func main() {
	// Serve static files (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialize the database
	db.Init()

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", middleware.RequireLogin(dashboardHandler))
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/notes", middleware.RequireLogin(notesHandler))
	http.HandleFunc("/notes/edit", middleware.RequireLogin(editNoteHandler))
	http.HandleFunc("/notes/delete", middleware.RequireLogin(deleteNoteHandler))

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
