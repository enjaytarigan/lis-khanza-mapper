package handler

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templateFS embed.FS

var baseTpl *template.Template

func initTemplates() error {
	t, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return err
	}
	baseTpl = t
	return nil
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := baseTpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

