package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"math/rand"
)

var dockerClient *client.Client
var containerPortMap = make(map[string][]*StalkerPort)
var containerMap = make(map[string]*StalkerContainer)
var tokenMap = make(map[string]string)

const DefaultTokenExpiry = 21600000

var EnvPassword = os.Getenv("PASSWORD")
var EnvPort = os.Getenv("PORT")
var EnvTokenExpiry = os.Getenv("TOKEN_EXPIRY_MILLI")
var AppBuildFolder = os.Getenv("APP_BUILD_FOLDER")

type Password struct {
	Password string `json:"password"`
}

type IsSecure struct {
	IsSecure bool `json:"isSecure"`
}

type PassToken struct {
	Token string `json:"token"`
}

func detailContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerId := vars["containerId"]

	fullContainerInfo, err := dockerClient.ContainerInspect(context.Background(), containerId)

	if err != nil {
		if client.IsErrContainerNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			panic(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}

	//get all network names
	networkNames := make([]string, len(fullContainerInfo.NetworkSettings.Networks))

	i := 0
	for k := range fullContainerInfo.NetworkSettings.Networks {
		networkNames[i] = k
		i++
	}

	ports, exists := containerPortMap[fullContainerInfo.ID]

	containerDetail := StalkerContainerDetail{
		Mounts:           GetStalkerMounts(fullContainerInfo.Mounts),
		Networks:         networkNames,
		EnvVars:          fullContainerInfo.Config.Env,
		StalkerContainer: containerMap[fullContainerInfo.ID],
	}

	if exists {
		containerDetail.Ports = ports
	} else {
		containerDetail.Ports = []*StalkerPort{}
	}

	json.NewEncoder(w).Encode(containerDetail)
}

func getAllContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var allContainers = []StalkerContainer{}

	for _, dc := range containers {
		containerPortMap[dc.ID] = GetStalkerPorts(dc.Ports)

		c := &StalkerContainer{
			Name:        dc.Names[0],
			Image:       dc.Image,
			Created:     dc.Created,
			Status:      dc.Status,
			ContainerId: dc.ID,
			State:       dc.State,
		}

		containerMap[dc.ID] = c
		allContainers = append(allContainers, *c);
	}

	json.NewEncoder(w).Encode(allContainers)
}

func tokenGenerator() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//POST call, takes in json body: { "password" : xxxxx }
//returns 200 and a 6-hour token for correct password
//returns 401 unauthorized otherwise
func login(w http.ResponseWriter, r *http.Request) {
	var password Password

	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if EnvPassword == password.Password {
		w.WriteHeader(http.StatusOK)
		token := tokenGenerator()
		tokenMap[token] = ""

		json.NewEncoder(w).Encode(&PassToken{Token: token})

		var tokenExpiryMilli = DefaultTokenExpiry

		if len(EnvTokenExpiry) > 0 {
			tokenExpiryMilli, err = strconv.Atoi(EnvTokenExpiry)
			if err != nil {
				tokenExpiryMilli = DefaultTokenExpiry
			}
		}

		tokenExpiryTimer := time.NewTimer(time.Millisecond * time.Duration(tokenExpiryMilli))
		go func() {
			<-tokenExpiryTimer.C
			delete(tokenMap, token)
		}()
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

func isSecure(w http.ResponseWriter, r *http.Request) {
	//GET call
	//returns true/false for whether or not there is a login
	//if PASSWORD is set, then true
	//returns { "isSecure" : false/true }

	if len(EnvPassword) > 0 {
		s := IsSecure{IsSecure: true}
		json.NewEncoder(w).Encode(s)
	} else {
		s := IsSecure{IsSecure: false}
		json.NewEncoder(w).Encode(s)
	}
}

func restartContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerId := vars["containerId"]

	log.Printf("Restarting container %s", containerId)

	waitDuration := 5 * time.Second
	err := dockerClient.ContainerRestart(context.Background(), containerId, &waitDuration)

	if err != nil {
		panic(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func ReturnJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSecure := len(EnvPassword) > 0
		token := r.Header.Get("Authorization")
		_, tokenExists := tokenMap[token]

		if !isSecure || tokenExists {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}

//serves compiled frontend code
func appHandler(w http.ResponseWriter, r *http.Request) {

	var buildPath = "../front/build"

	if len(AppBuildFolder) > 0 {
		buildPath = AppBuildFolder
	}

	var fullBuildPath = buildPath + "/index.html"

	log.Printf("Trying to serve root index file at: %s", fullBuildPath)

	http.ServeFile(w, r, fullBuildPath)
}

func main() {
	//initialize docker client
	dockerClient, _ = client.NewEnvClient()

	r := mux.NewRouter()


	var buildPath = "../front/build"

	if len(AppBuildFolder) > 0 {
		buildPath = AppBuildFolder
	}

	var fullBuildPath = buildPath + "/static"

	log.Printf("static files is located at: %s", fullBuildPath)


	r.HandleFunc("/", appHandler)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(buildPath + "/static"))))

	r.Handle("/api/containers",
		alice.New(Protected, ReturnJSON).Then(http.HandlerFunc(getAllContainers)))
	r.Handle("/api/container/{containerId}/restart",
		alice.New(Protected).Then(http.HandlerFunc(restartContainer))).Methods("POST")
	r.Handle("/api/container/{containerId}/detail",
		alice.New(Protected, ReturnJSON).Then(http.HandlerFunc(detailContainer)))

	r.Handle("/api/login", alice.New(ReturnJSON).Then(http.HandlerFunc(login))).Methods("POST", "OPTIONS")
	r.Handle("/api/isSecure", alice.New(ReturnJSON).Then(http.HandlerFunc(isSecure)))

	if len(EnvPort) == 0 {
		EnvPort = "8080"
	}

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Printf("Starting HTTP server on port %s", EnvPort)

	if err := http.ListenAndServe(":"+EnvPort, handlers.CORS(originsOk, methodsOk, headersOk)(r)); err != nil {
		panic(err)
	}
}
