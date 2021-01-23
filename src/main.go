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
	"regexp"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func getDir() string {
	rootDir := os.Getenv("API_ROOT_DIRECTORY")
	if rootDir == "" {
		rootDir = "./images/"
	}
	return rootDir
}

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
	rootDir := getDir()
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

func retrieveFilePaths(w http.ResponseWriter, r *http.Request) {
	protocol := ""
	if r.TLS == nil {
		protocol = "http://"
	} else {
		protocol = "https://"
	}

	params := mux.Vars(r)
	id := params["id"]

	rootDir := getDir()
	dir := rootDir + id
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

	return http.StripPrefix("/v1/files/", http.FileServer(http.Dir(getDir())))
}

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("listening on", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/files", uploadFile).
		Methods("POST")
	router.HandleFunc("/v1/files/{id:[A-z0-9-]+}", retrieveFilePaths).
		Methods("GET")
	router.PathPrefix("/v1/files/").Handler(retrieveFile()).Methods("GET")
	http.ListenAndServe(":"+port, router)
}
