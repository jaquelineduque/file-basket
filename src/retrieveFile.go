package main

import "net/http"

func retrieveFile() http.Handler {

	return http.StripPrefix("/v1/files/", http.FileServer(http.Dir(GetDir())))
}
