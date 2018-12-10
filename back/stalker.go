package main

import (
	"net/http"

	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"encoding/json"
	"github.com/gorilla/mux"
)

var dockerClient *client.Client
var containerMap = make(map[string][]*StalkerPort)

func detailContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	containerDetail := StalkerContainerDetail {
		Ports: containerMap[containerId],
		Mounts: GetStalkerMounts(fullContainerInfo.Mounts),
		ContainerId: containerId,
		Networks: networkNames,
		EnvVars: fullContainerInfo.Config.Env,
	}


	json.NewEncoder(w).Encode(containerDetail)
}

func getAllContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	//initialize docker client
	dockerClient, _ = client.NewEnvClient()

	//TODO: set access control origin for the backend

	//TODO: build something that checks for password on all endpoints
	http.HandleFunc("/containers", getAllContainers)
	http.HandleFunc("/container/{containerId}/restart", restartContainer)
	http.HandleFunc("/container/{containerId}/detail", detailContainer)

	http.HandleFunc("/login", login)
	http.HandleFunc("/isSecure", isSecure)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}