package main

import (
	"fmt"
	"go-web-app/db"
	"go-web-app/handlers"
	"go-web-app/helper"
	"go-web-app/middleware"

	"net/http"

	"github.com/gorilla/sessions"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	helper.RenderTemplate(w, "home.html", nil)
}

func main() {
	// Serve static files (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialize the database
	db.Init()

	// Init session store
	store := sessions.NewCookieStore([]byte("super-secret-key"))
	handlers.InitStore(store)
	middleware.InitStore(store)

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/dashboard", middleware.RequireLogin(handlers.DashboardHandler))
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/notes", middleware.RequireLogin(handlers.NotesHandler))
	http.HandleFunc("/notes/edit", middleware.RequireLogin(handlers.EditNoteHandler))
	http.HandleFunc("/notes/delete", middleware.RequireLogin(handlers.DeleteNoteHandler))

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
