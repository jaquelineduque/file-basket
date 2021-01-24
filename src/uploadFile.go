package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	if id != "" {
		r, _ := regexp.Compile("^([A-z0-9-]+)$")
		if !r.MatchString(id) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Id should contains only letters, numbers and hyphen"))
			return
		}
	}
	rootDir := GetDir()
	diretorio := rootDir + id

	err = os.MkdirAll(diretorio, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory")
		fmt.Println(err)
		return
	}

	defer file.Close()

	completeFilePath := diretorio + "/" + handler.Filename
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	destinyPath, err := os.Create(completeFilePath)
	defer destinyPath.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(destinyPath, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File on\n")
}
