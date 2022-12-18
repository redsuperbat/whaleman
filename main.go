package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
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
	DATA_DIR = "/var/lib/whaleman"
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

func toMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func toFilename(url string) string {
	return DATA_DIR + "/" + toMD5Hash(url)
}

func readCache(url string) []byte {
	filepath := toFilename(url)
	b, err := ioutil.ReadFile(filepath)

	if os.IsNotExist(err) {
		log.Println("File did not exist. Creating an empty one to start with")
		log.Println("Ensuring", DATA_DIR, "exists.")
		if err = os.MkdirAll(DATA_DIR, 0700); err != nil {
			log.Fatalln(err)
		}
		initBytes := []byte("")
		log.Println("Creating empty file", filepath)
		err := ioutil.WriteFile(filepath, initBytes, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		return initBytes
	}

	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func checkFile(url string) error {
	b := downloadGithubFile(url)
	b2 := readCache(url)
	if filesEqual(&b, &b2) {
		log.Println("Remote files match local cache.")
		return nil
	}
	log.Println("Mismatch against local cache. Updating cache.")
	filepath := toFilename(url)
	if err := ioutil.WriteFile(filepath, b, 0644); err != nil {
		return err
	}

	log.Println("Trying to restart docker applications")
	cmd := exec.Command("docker-compose", "-f", filepath, "up", "-d")
	log.Println("Running command with args: ", cmd.Args)
	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println(m)
	}
	cmd.Wait()
	return nil
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
			err := checkFile(u)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}()
	}
}

func main() {
	port := ":8090"
	log.Println("Registered handler for path {/}")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to check files")
		checkFiles()
		io.WriteString(w, "Thanks for your request")
	})
	log.Println("Server started on port", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
