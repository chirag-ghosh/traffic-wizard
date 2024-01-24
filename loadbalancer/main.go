package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	// "path/filepath"
	"github.com/chirag-ghosh/traffic-wizard/loadbalancer/internal/consistenthashmap"
)

func spawnNewServerInstance(hostname string, id int) {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Current working directory:", dir)

	cmd := exec.Command("sudo", "docker", "build", "--tag", "traffic-wizard-server", "/server")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to build server image: %v", err)
	}

	// Run the server Docker container
	cmd = exec.Command("sudo", "docker", "run", "-d", "--name", hostname, "-e", fmt.Sprintf("ID=%d", id), "traffic-wizard-server:latest")

	// cmd = exec.Command("docker", "run", "-d", "--name", hostname, "traffic-wizard-server:latest")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start new server instance: %v", err)
	}
}

func hashRequest(path string) int {
	h := fnv.New32a()
	_, err := h.Write([]byte(path))
	if err != nil {
		return 0
	}
	return int(h.Sum32())
}

type ServerInfo struct {
	ID       int
	Hostname string
}

var chm = consistenthashmap.ConsistentHashMap{}
var servers = make(map[int]ServerInfo)

func init() {
	chm.Init()
	// Initialize your server instances here with actual IDs and hostnames
	// servers[1] = ServerInfo{ID: 1, Hostname: "server-1-hostname"}
	// ... add other servers
}

func getReplicaStatus(w http.ResponseWriter, r *http.Request) {
	replicas := make([]string, 0, len(servers))
	for _, serverInfo := range servers {
		replicas = append(replicas, serverInfo.Hostname)
	}

	response := map[string]interface{}{
		"message": map[string]interface{}{
			"N":        len(servers),
			"replicas": replicas,
		},
		"status": "successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getNextServerID() int {
	maxID := 0
	for id := range servers {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}

type AddServersPayload struct {
	N         int      `json:"n"`
	Hostnames []string `json:"hostnames"`
}

type AddServersResponse struct {
	Message map[string]interface{} `json:"message"`
	Status  string                 `json:"status"`
}

func getReplicas() []string {
	replicas := make([]string, 0, len(servers))
	for _, serverInfo := range servers {
		replicas = append(replicas, serverInfo.Hostname)
	}
	return replicas
}

func addServersEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	var payload AddServersPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform sanity checks on the request payload
	if len(payload.Hostnames) > payload.N {
		resp := AddServersResponse{
			Message: map[string]interface{}{"<Error>": "Length of hostname list is more than newly added instances"},
			Status:  "failure",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for i, hostname := range payload.Hostnames {
		serverID := getNextServerID()
		chm.AddServer(serverID)

		// The logic to actually spawn the server instances should be here.
		spawnNewServerInstance(hostname, serverID)

		servers[serverID] = ServerInfo{ID: serverID, Hostname: hostname}

		if i+1 == payload.N {
			break
		}
	}

	resp := AddServersResponse{
		Message: map[string]interface{}{
			"N":        len(servers),
			"replicas": getReplicas(),
		},
		Status: "successful",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type RemoveServersPayload struct {
	N         int      `json:"n"`
	Hostnames []string `json:"hostnames"`
}

type RemoveServersResponse struct {
	Message map[string]interface{} `json:"message"`
	Status  string                 `json:"status"`
}

func chooseRandomServer() string {
	rand.Seed(time.Now().UnixNano())

	keys := make([]int, 0, len(servers))
	for key := range servers {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return ""
	}

	randomServerID := keys[rand.Intn(len(keys))]

	return servers[randomServerID].Hostname
}

func removeServerInstance(hostname string) {
    cmd := exec.Command("sudo", "docker", "stop", hostname)
    err := cmd.Run()
    if err != nil {
        log.Fatalf("Failed to stop server instance '%s': %v", hostname, err)
    }

    cmd = exec.Command("sudo", "docker", "rm", hostname)
    err = cmd.Run()
    if err != nil {
        log.Fatalf("Failed to remove server instance '%s': %v", hostname, err)
    }

    var serverID int
    for id, info := range servers {
        if info.Hostname == hostname {
            serverID = id
            break
        }
    }
    delete(servers, serverID)
    chm.RemoveServer(serverID)
}
func removeServersEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var payload RemoveServersPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(payload.Hostnames) > payload.N {
		resp := RemoveServersResponse{
			Message: map[string]interface{}{"<Error>": "Length of hostname list is more than removable instances"},
			Status:  "failure",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for _, hostname := range payload.Hostnames {
		removeServerInstance(hostname)
	}

	for len(payload.Hostnames) < payload.N {
		hostname := chooseRandomServer()
		removeServerInstance(hostname)
		payload.Hostnames = append(payload.Hostnames, hostname)
	}

	resp := RemoveServersResponse{
		Message: map[string]interface{}{
			"N":        len(servers),
			"replicas": getReplicas(),
		},
		Status: "successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func responseError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
		"status":  "failure",
	})
}

// todo

func routeRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	requestID := hashRequest(path)
	serverID := chm.GetServerForRequest(requestID)

	server, exists := servers[serverID]
	if !exists {
		responseError(w, "<Error> Server not found", http.StatusNotFound)
		return
	}

	resp, err := http.Get("http://" + server.Hostname + path)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "<Error> '"+path+"' endpoint does not exist in server replicas", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	http.HandleFunc("/rep", getReplicaStatus)
	http.HandleFunc("/add", addServersEndpoint)
	http.HandleFunc("/rm", removeServersEndpoint)
	http.HandleFunc("/", routeRequest)

	fmt.Println("Load Balancer started on port 5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}
