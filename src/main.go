package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
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
	diretorio := "images/" + id

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

func retrieveFilePaths(w http.ResponseWriter, r *http.Request) {
	protocol := ""
	if r.TLS == nil {
		protocol = "http://"
	} else {
		protocol = "https://"
	}

	params := mux.Vars(r)
	id := params["id"]
	dir := "./images/" + id
	files, _ := ioutil.ReadDir(dir)

	var arquivos []Path

	if len(files) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Id doesn't have files"))
		return
	}

	for _, file := range files {
		var path Path
		path.Path = protocol + filepath.Join(r.Host, r.URL.Path, file.Name())
		arquivos = append(arquivos, path)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(arquivos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Fail returning paths"))
	}

}

func retrieveFile() http.Handler {
	return http.StripPrefix("/v1/files/", http.FileServer(http.Dir("./images")))
}

func main() {
	port := ":8080"
	log.Println("listening on", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/files", uploadFile).
		Methods("POST")
	router.HandleFunc("/v1/files/{id}", retrieveFilePaths).
		Methods("GET")
	router.PathPrefix("/v1/files/").Handler(retrieveFile()).Methods("GET")
	http.ListenAndServe(port, router)
}
