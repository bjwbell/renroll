package main
import (
	"net/http"
	"html/template"
)

type Index struct {
	GoogleAnalyticsId string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := Index{GoogleAnalyticsId: configuration().GoogleAnalyticsId}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, index)
}
