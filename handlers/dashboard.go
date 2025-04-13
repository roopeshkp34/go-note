package handlers

import (
	"go-web-app/db"
	"go-web-app/helper"
	"go-web-app/models"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	db.DB.Find(&notes)
	helper.RenderTemplate(w, "dashboard.html", notes)
}
