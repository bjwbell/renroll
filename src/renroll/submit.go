package renroll

import (
	"html/template"
	"net/http"
)

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(
		"submit.html",
		"templates/header.html",
		"templates/bottombar.html")
	t.Execute(w, struct{ Conf Configuration }{Config()})
}
