package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var rv = ResumeValidator{}
var jsonTemplate = `
{
  "basics": {
    "name": "John Doe",
    "label": "Programmer",
    "email": "john@gmail.com"
  }
}
`

func TestResumeValidator_HomePage(t *testing.T) {
	rv.RegisterValidatorAPI()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("err not expected: %v", err)
	}

	handler := rv.homePage()
	assertServiceCall(t, handler, req, http.StatusOK)
}

func TestResumeValidator_Validatefile_RedirectToHomePage(t *testing.T) {
	req, err := http.NewRequest("POST", "/validateFile", nil)
	if err != nil {
		t.Fatalf("err not expected: %v", err)
	}

	handler := rv.validateFile()
	assertServiceCall(t, handler, req, http.StatusTemporaryRedirect)
}

func TestResumeValidator_ValidateFile_ValidResume(t *testing.T) {
	f, cleanup := mustCreateTempFile(t, jsonTemplate)
	defer cleanup(f)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		t.Fatalf("could not create form file: %v", err)
	}

	if _, err = io.Copy(fw, f); err != nil {
		t.Fatalf("could not copy file into buffer: %v", err)
	}
	w.Close()

	req, err := http.NewRequest("POST", "/validateFile", &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	handler := rv.validateFile()
	assertServiceCall(t, handler, req, http.StatusOK)
}

func TestResumeValidator_ValidateFile_InvalidResume(t *testing.T) {
	f, cleanup := mustCreateTempFile(t, "{}}")
	defer cleanup(f)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		t.Fatalf("could not create form file: %v", err)
	}

	if _, err = io.Copy(fw, f); err != nil {
		t.Fatalf("could not copy file into buffer: %v", err)
	}
	w.Close()

	req, err := http.NewRequest("POST", "/validateFile", &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	handler := rv.validateFile()
	assertServiceCall(t, handler, req, http.StatusOK)
}

func mustCreateTempFile(t *testing.T, body string) (*os.File, func(file *os.File)) {
	content := []byte(body)
	tmp, err := ioutil.TempFile("", "testing_tmp")
	if err != nil {
		t.Fatalf("err creating temp file: %v", err)
	}

	if _, err := tmp.Write(content); err != nil {
		t.Fatalf("err writing into the temp file: %v", err)
	}

	return tmp, func(file *os.File) {
		file.Close()
		os.Remove(file.Name())
	}
}

func assertServiceCall(t *testing.T, handler http.HandlerFunc, r *http.Request, expStatusCode int) {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	if status := rr.Code; status != expStatusCode {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, expStatusCode)
	}
}
