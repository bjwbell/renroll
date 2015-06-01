package main
import (
	"net/http"
	"html/template"
)

type Submit struct {
	Conf Configuration
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		"submit.html",
		"templates/header.html",
		"templates/bottombar.html")
	t.Execute(w, Submit{Conf: configuration()})
}
