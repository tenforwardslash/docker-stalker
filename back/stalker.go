package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"os"
)

type Password struct {
	Password string `json:"password"`
}

func getAllContainers(w http.ResponseWriter, r *http.Request) {
	//returns a list of docker containers

}

func login(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPost) {
		return
	}

	//POST call
	//this takes in a json body: { "password" : xxxxx }
	//the passed body is compared against environment variable PASSWORD set on backend startup
	//returns 200 for correct password, 401 unauthorized

	var password Password

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		panic(err)
	}

	fmt.Printf("you posted password: %s", password.Password)

	PASSWORD := os.Getenv("PASSWORD")

	if PASSWORD == password.Password {
		w.WriteHeader(200)
		w.Write([]byte("200 - all good"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("401 - unathorized"))
	}

}

func isSecure(w http.ResponseWriter, r *http.Request) {
	//GET call
	//returns true/false for whether or not there is a login
	//if PASSWORD is set, then true
	//returns { "isSecure" : false/true }
}

func restartContainer(w http.ResponseWriter, r *http.Request) {
	//POST call
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/containers", getAllContainers)
	r.HandleFunc("/restart", restartContainer)
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/isSecure", isSecure)
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}