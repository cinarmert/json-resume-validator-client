package main

import (
	"bytes"
	resumeValidator "github.com/cinarmert/json-resume-validator"
	"io"
	"log"
	"net/http"
)

// ResumeValidator is the a controller type
type ResumeValidator struct{}

// RegisterValidatorAPI registers the endpoints of ResumeValidator controller
func (rv *ResumeValidator) RegisterValidatorAPI() {
	http.HandleFunc("/", rv.homePage())
	http.HandleFunc("/validateFile", rv.validateFile())
}

func (rv *ResumeValidator) homePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "Result will be shown here!"
		render(w, "home", Page{Status: status})
	}
}

func (rv *ResumeValidator) validateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("could not get file, redirecting...")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		var buf bytes.Buffer
		defer file.Close()
		io.Copy(&buf, file)

		rv := new(resumeValidator.ResumeValidator).WithData(buf.Bytes())
		status := "Not a valid JSON Resume!"
		if rv.IsValid() {
			status = "Valid JSON Resume!"
		}

		render(w, "home", Page{Status: status})
	}
}
