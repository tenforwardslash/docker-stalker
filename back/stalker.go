package main

import (
	"net/http"

	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"encoding/json"
	"fmt"
)

func getAllContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//returns a list of docker containers
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), 	types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var allContainers = []StalkerContainer{}

	for _, dc := range containers {
		//TODO: get the environment variables
		//TODO: get the network
		fmt.Printf("%s", dc.Names)
		c := StalkerContainer {
			Name: dc.Names[0],
			Image: dc.Image,
			Created: dc.Created,
			Status: dc.Status,
			Ports: GetStalkerPorts(dc.Ports),
			Mounts: GetStalkerMounts(dc.Mounts),
			ContainerId: dc.ID,
			State: dc.State,
			//Network:
			//EnvVars:
		}

		allContainers = append(allContainers, c);
		//fmt.Printf("%s %s\n", container.ID[:10], container.Image)
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
	http.HandleFunc("/containers", getAllContainers)
	http.HandleFunc("/restart", restartContainer)
	http.HandleFunc("/login", login)
	http.HandleFunc("/isSecure", isSecure)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}