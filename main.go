package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Page holds the parameters to the html file
type Page struct {
	Status string
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	addr := listenAddr + `:` + os.Getenv("PORT")

	validator := ResumeValidator{}
	validator.RegisterValidatorAPI()

	log.Printf("starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func render(w http.ResponseWriter, tmpl string, page Page) {
	tmpl = fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFiles(tmpl)

	if err != nil {
		log.Print("template parsing error: ", err)
		return
	}

	_ = t.Execute(w, page)
}
