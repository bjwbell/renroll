package main
import (
	"net/http"
	"html/template"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("submit.html")
	t.Execute(w, configuration())
}
