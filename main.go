package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	DATA_DIR     = "/var/lib/whaleman"
	COMPOSE_FILE = DATA_DIR + "/compose-file"
)

func downloadGithubFile(url string) []byte {
	ghToken := os.Getenv("GH_PAT")
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "token "+ghToken)
	request.Header.Add("Accept", "application/vnd.github.v3.raw")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

func checkFile(url string) {
	b := downloadGithubFile(url)
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

func getUrls() []string {
	urls := strings.Split(os.Getenv("GH_COMPOSE_FILES"), ",")
	if len(urls) == 0 {
		log.Fatalln("Please supply the env variable GH_COMPOSE_FILES")
	}
	return urls
}

func main() {
	http.HandleFunc("/handle-changes", func(w http.ResponseWriter, r *http.Request) {
		urls := getUrls()
		var wg sync.WaitGroup
		wg.Add(len(urls))

		for _, url := range urls {
			u := url
			// Run every url request in parallell
			go func() {
				checkFile(u)
				wg.Done()
			}()
		}
	})
	log.Println("Started server")
	http.ListenAndServe(":8090", nil)
}
