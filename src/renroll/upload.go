package renroll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type UploadTemplate struct {
	Conf Configuration
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("uploadHandler - begin")
	conf := Config()
	conf.GPlusSigninCallback = "gDummy"
	conf.FacebookSigninCallback = "fbDummy"
	upload := UploadTemplate{conf}
	t, _ := template.ParseFiles(
		"upload.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	log.Print("uploadHandler - execute")

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err == nil {
		fmt.Println("UPLOADDDDDDDDD")
		fmt.Println("file:", file)
		fmt.Println("header.Filename:", header.Filename)
		NotifyUpload(header.Filename)
	} else {
		t.Execute(w, upload)
	}
}

func NotifyUpload(filename string) {
	userEmail := "unknown@unknown.com"
	subject := "Renroll Lease Uploaded (" + filename + ")"
	body := "Lease filename: " + filename + ".\r\n"
	SendAdminEmail(userEmail, subject, body)
}
