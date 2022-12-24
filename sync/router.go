package sync

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/redsuperbat/whaleman/data"
)

type Msg struct {
	Message string `json:"message"`
}

func downloadGithubFile(url string) ([]byte, error) {
	ghToken := os.Getenv("GH_PAT")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "token "+ghToken)
	request.Header.Add("Accept", "application/vnd.github.v3.raw")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func filesEqual(b *[]byte, b2 *[]byte) bool {
	return bytes.Equal(*b, *b2)
}

func toMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func sumChars(str string) int64 {
	sum := 0
	for _, char := range str {
		sum += int(char)
	}
	return int64(sum)
}

func toFilename(url string) string {
	return toMD5Hash(url)
}

func readCache(url string) []byte {
	filename := toFilename(url)
	if err, b := data.ReadManifestFile(filename); err != nil {
		return []byte("")
	} else {
		return b
	}
}

func checkFile(log *golog.Logger, url string) error {
	b, err := downloadGithubFile(url)
	if err != nil {
		return err
	}
	b2 := readCache(url)
	if filesEqual(&b, &b2) {
		log.Info("Remote files match local cache.")
		return nil
	}

	log.Info("Mismatch against local cache. Updating cache.")
	filepath := toFilename(url)
	if err = data.WriteManifestFile(filepath, b); err != nil {
		return err
	}

	log.Info("Trying to restart docker applications")
	seed := sumChars(toMD5Hash(url))
	faker := gofakeit.New(seed)
	project := strings.ToLower(faker.Adjective() + "-" + faker.Animal())
	log.Info("Project", project, "generated")
	cmd := exec.Command("docker-compose", "-f", filepath, "-p", project, "up", "-d")
	log.Info("Running command with args: ", cmd.Args)
	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println(m)
	}
	cmd.Wait()

	if cmd.ProcessState.ExitCode() == 0 {
		log.Println("Updated docker containers.")
		return nil
	}

	if err := os.Remove(filepath); err != nil {
		log.Info("Unable to restart docker containers with new manifest. Invalidating cache.")
		return err
	}
	return errors.New("Unable to update docker containers")
}

func getUrls(log *golog.Logger) (error, []string) {
	if err, urls := data.ReadManifestResources(); err != nil {
		return err, nil
	} else {
		log.Println("Urls: ", urls)

		return nil, urls
	}
}

func checkFiles(log *golog.Logger) {
	err, urls := getUrls(log)

	if err != nil {
		log.Error(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		u := strings.TrimSpace(url)
		// Run every url request in parallell
		log.Info("Checking file", url)
		go func() {
			err := checkFile(log, u)
			if err != nil {
				log.Error(err)
			}
			wg.Done()
		}()
	}
}

func startPoll(log *golog.Logger) {
	pollInterval := os.Getenv("POLLING_INTERVAL_MIN")
	if pollInterval == "" {
		log.Debug("Polling disabled")
		return
	}

	log.Debug("Polling enabled polling every", pollInterval, "minutes")
	interval, err := strconv.Atoi(pollInterval)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	for {
		select {
		case <-ticker.C:
			checkFiles(log)
		}
	}
}

func handleSync(ctx iris.Context) {
	checkFiles(ctx.Application().Logger())
	ctx.JSON(Msg{Message: "Successfully synced!"})
}

func RegisterSync(app *iris.Application) {
	syncApi := app.Party("/sync")
	syncApi.Post("/", handleSync)
	go startPoll(app.Logger())
}
