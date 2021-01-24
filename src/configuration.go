package main

import "os"

func GetDir() string {
	rootDir := os.Getenv("API_ROOT_DIRECTORY")
	if rootDir == "" {
		rootDir = "./images/"
	}
	return rootDir
}

func GetPort() string {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
