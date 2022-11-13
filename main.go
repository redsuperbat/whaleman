package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	DATA_DIR     = "/var/lib/whaleman"
	COMPOSE_FILE = DATA_DIR + "/compose-file"
)

func handleComposeChanges(writer http.ResponseWriter, req *http.Request) {
	fileUrl := os.Getenv("COMPOSE_URL")
	if fileUrl == "" {
		log.Fatalln("Please supply the env variable COMPOSE_URL")
	}
	response, err := http.Get(fileUrl)

	if err != nil {
		log.Fatalln(err)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	b2, err := ioutil.ReadFile(COMPOSE_FILE)
	if err != nil {
		log.Fatalln(err)
	}

	if bytes.Equal(b2, b) {
		return
	}

	if err = ioutil.WriteFile(COMPOSE_FILE, b, 0644); err != nil {
		log.Fatalln(err)
	}

	exec.Command("docker-compose", "restart", "-f", COMPOSE_FILE)
}

func main() {
	http.HandleFunc("/handle-changes", handleComposeChanges)
	log.Println("Started server")
	http.ListenAndServe(":8090", nil)
}
