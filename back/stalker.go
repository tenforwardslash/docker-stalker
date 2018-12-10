package main

import (
	"net/http"
	"os"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"encoding/json"
	"github.com/gorilla/mux"
	"time"
	"fmt"
)

var dockerClient *client.Client
var containerMap = make(map[string][]*StalkerPort)

func detailContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerId := vars["containerId"]

	fullContainerInfo, err := dockerClient.ContainerInspect(context.Background(), containerId)

	if err != nil {
		panic(err)
	}

	//get all network names
	networkNames := make([]string, len(fullContainerInfo.NetworkSettings.Networks))

	i := 0
	for k := range fullContainerInfo.NetworkSettings.Networks {
		networkNames[i] = k
		i++
	}

	ports, exists := containerMap[containerId]

	containerDetail := StalkerContainerDetail {
		Mounts: GetStalkerMounts(fullContainerInfo.Mounts),
		ContainerId: containerId,
		Networks: networkNames,
		EnvVars: fullContainerInfo.Config.Env,
	}

	if exists {
		containerDetail.Ports = ports
	} else {
		containerDetail.Ports = []*StalkerPort{}
	}


	json.NewEncoder(w).Encode(containerDetail)
}


type Password struct {
	Password string `json:"password"`
}


func getAllContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := dockerClient.ContainerList(context.Background(), 	types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var allContainers = []StalkerContainer{}

	for _, dc := range containers {
		containerMap[dc.ID] = GetStalkerPorts(dc.Ports)

		c := StalkerContainer {
			Name: dc.Names[0],
			Image: dc.Image,
			Created: dc.Created,
			Status: dc.Status,
			ContainerId: dc.ID,
			State: dc.State,
		}

		allContainers = append(allContainers, c);
	}

	json.NewEncoder(w).Encode(allContainers)
}

func login(w http.ResponseWriter, r *http.Request) {
	//POST call
	//this takes in a json body: { "password" : xxxxx }
	//the passed body is compared against environment variable PASSWORD set on backend startup
	//returns 200 for correct password, 401 unauthorized
	if (r.Method != http.MethodPost) {
		return
	}

	var password Password

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("you posted password: %s", password.Password)

	PASSWORD := os.Getenv("PASSWORD")

	if PASSWORD == password.Password {
		w.WriteHeader(200)
		w.Write([]byte("200 - all good"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("401 - unathorized"))
	}

}

type IsSecure struct {
	IsSecure bool `json:"isSecure"`
}

func isSecure(w http.ResponseWriter, r *http.Request) {
	//GET call
	//returns true/false for whether or not there is a login
	//if PASSWORD is set, then true
	//returns { "isSecure" : false/true }

	PASSWORD := os.Getenv("PASSWORD")

	if len(PASSWORD) != 0 {
		s := IsSecure{ IsSecure: true }
		json.NewEncoder(w).Encode(s)
	} else {
		s := IsSecure{ IsSecure: false }
		json.NewEncoder(w).Encode(s)
	}
}

func restartContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerId := vars["containerId"]

	fmt.Printf("inside of restart container %s", containerId)

	waitDuration := 5 * time.Second
	err := dockerClient.ContainerRestart(context.Background(), containerId, &waitDuration)

	if err != nil {
		panic(err)
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
	}

}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	//initialize docker client
	dockerClient, _ = client.NewEnvClient()

	r := mux.NewRouter()

	//TODO: build something that checks for password on all endpoints
	r.HandleFunc("/containers", getAllContainers)
	r.HandleFunc("/container/{containerId}/restart", restartContainer).Methods("POST")
	r.HandleFunc("/container/{containerId}/detail", detailContainer)

	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/isSecure", isSecure)

	r.Use(CORS)
	
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}