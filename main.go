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

	"github.com/joho/godotenv"
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

func filesEqual(b *[]byte, b2 *[]byte) bool {
	return bytes.Equal(*b, *b2)
}

func removeEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func checkFile(url string) {
	b := downloadGithubFile(url)
	b2, err := ioutil.ReadFile(COMPOSE_FILE)
	if err != nil {
		log.Fatalln(err)
	}

	if filesEqual(&b, &b2) {
		log.Println("Remote files match local cache.")
		return
	}
	log.Println("Mismatch against local cache. Updating cache.")
	if err = ioutil.WriteFile(COMPOSE_FILE, b, 0644); err != nil {
		log.Fatalln(err)
	}

	log.Println("Executing command to restart affected application")
	exec.Command("docker-compose", "restart", "-f", COMPOSE_FILE)
}

func getUrls() []string {
	urls := removeEmpty(strings.Split(os.Getenv("GH_COMPOSE_FILES"), ","))
	log.Println("Urls: ", urls)
	if len(urls) == 0 {
		log.Fatalln("Please supply the env variable GH_COMPOSE_FILES")
	}
	return urls
}

func checkFiles() {
	urls := getUrls()
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		u := url
		// Run every url request in parallell
		log.Println("Checking file", url)
		go func() {
			checkFile(u)
			wg.Done()
		}()
	}
}

func main() {
	port := ":8090"
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file:", err)
	}
	log.Println("Registered handler for path {/}")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to check files")
		checkFiles()
		io.WriteString(w, "Thanks for your request")
	})
	log.Println("Server started on port", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
