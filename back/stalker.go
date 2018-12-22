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
	"github.com/justinas/alice"
	"github.com/urfave/negroni"
	"crypto/rand"
	"strconv"
	"io/ioutil"
)

var dockerClient *client.Client
var containerPortMap = make(map[string][]*StalkerPort)
var containerMap = make(map[string]*StalkerContainer)
var tokenMap = make(map[string]string)

const DefaultTokenExpiry = 21600000

type Password struct {
	Password string `json:"password"`
}

type IsSecure struct {
	IsSecure bool `json:"isSecure"`
}

type PassToken struct {
	Token string `json:"token"`
}

//returns detail for specified container in path
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

	ports, exists := containerPortMap[containerId]

	containerDetail := StalkerContainerDetail {
		Mounts: GetStalkerMounts(fullContainerInfo.Mounts),
		Networks: networkNames,
		EnvVars: fullContainerInfo.Config.Env,
		StalkerContainer: containerMap[containerId],
	}

	if exists {
		containerDetail.Ports = ports
	} else {
		containerDetail.Ports = []*StalkerPort{}
	}


	json.NewEncoder(w).Encode(containerDetail)
}

//returns list of all running docker containers
func getAllContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := dockerClient.ContainerList(context.Background(), 	types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var allContainers = []StalkerContainer{}

	for _, dc := range containers {
		containerPortMap[dc.ID] = GetStalkerPorts(dc.Ports)

		c := &StalkerContainer {
			Name: dc.Names[0],
			Image: dc.Image,
			Created: dc.Created,
			Status: dc.Status,
			ContainerId: dc.ID,
			State: dc.State,
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
		return
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &password)
		if err != nil {
			panic(err)
			return
		}
	}

	fmt.Printf("sdjfklsdf %s", password.Password)

	PASSWORD := os.Getenv("PASSWORD")

	if PASSWORD == password.Password {
		w.WriteHeader(http.StatusOK)
		token := tokenGenerator()
		tokenMap[token] = ""

		json.NewEncoder(w).Encode(&PassToken{Token: token})

		var tokenExpiryMilli = DefaultTokenExpiry
		tokenExpiryEnv := os.Getenv("TOKEN_EXPIRY_MILLI")

		if len(tokenExpiryEnv) > 0 {
			tokenExpiryMilli, err = strconv.Atoi(tokenExpiryEnv)
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

//returns { "isSecure" : false/true } on whether or not backend is locked down
func isSecure(w http.ResponseWriter, r *http.Request) {
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

func ReturnJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("PASSWORD")
		isSecure := len(pass) > 0
		token := w.Header().Get("Authorization")
		_, tokenExists := tokenMap[token]

		if !isSecure || tokenExists {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}

func main() {
	//initialize docker client
	dockerClient, _ = client.NewEnvClient()

	r := mux.NewRouter()

	r.Use(CORS)

	r.Handle("/containers",
		alice.New(Protected, ReturnJSON).Then(http.HandlerFunc(getAllContainers)))
	r.Handle("/container/{containerId}/restart",
		alice.New(Protected).Then(http.HandlerFunc(restartContainer))).Methods("POST")
	r.Handle("/container/{containerId}/detail",
		alice.New(Protected, ReturnJSON).Then(http.HandlerFunc(detailContainer)))

	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/isSecure", isSecure)


	n := negroni.Classic()
	n.UseHandler(r)

	n.Run(":8080")
}