package handlers

import (
	"go-web-app/db"
	"go-web-app/models"
	"net/http"
)

func NotesHandler(w http.ResponseWriter, r *http.Request) {
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

func EditNoteHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	db.DB.Delete(&models.Note{}, id)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
