package main

import (
"net/http"
"strings"
)

func getAllContainers(w http.ResponseWriter, r *http.Request) {
	//returns a list of docker containers
}

func login(w http.ResponseWriter, r *http.Request) {
	//POST call
	//this takes in a json body: { "password" : xxxxx }
	//the passed body is compared against environment variable PASSWORD set on backend startup
	//returns true/false auth
}

func isSecure(w http.ResponseWriter, r *http.Request) {
	//GET call
	//returns true/false for whether or not there is a login
	//if PASSWORD is set, then true
}


func restartContainer(w http.ResponseWriter, r *http.Request) {
	//POST call
}

func main() {
	http.HandleFunc("/containers", getAllContainers)
	http.HandleFunc("/restart", restartContainer)
	http.HandleFunc("/login", login)
	http.HandleFunc("/isSecure", isSecure)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}