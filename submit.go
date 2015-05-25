package main
import (
	"net/http"
	"html/template"
)

type SubmitData struct {
	GoogleAnalyticsId string
	GoogleClientID string
}
func submitHandler(w http.ResponseWriter, r *http.Request) {
	conf := configuration()
	
	submitData := SubmitData{GoogleAnalyticsId: conf.GoogleAnalyticsId,
		GoogleClientID: conf.GoogleClientID}
	
	t, _ := template.ParseFiles("submit.html")
	t.Execute(w, submitData)
}
