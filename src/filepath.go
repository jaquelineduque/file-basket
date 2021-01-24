package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func retrieveFilePaths(w http.ResponseWriter, r *http.Request) {
	protocol := ""
	if r.TLS == nil {
		protocol = "http://"
	} else {
		protocol = "https://"
	}

	params := mux.Vars(r)
	id := params["id"]

	rootDir := GetDir()
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
