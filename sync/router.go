package sync

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redsuperbat/whaleman/data"
	"github.com/redsuperbat/whaleman/docker"
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

func readCache(url string) ([]byte, error) {
	filename := toFilename(url)
	return data.ReadManifestFile(filename)
}

func getProjectNameFromHash(hash string) string {
	seed := sumChars(hash)
	faker := gofakeit.New(seed)
	return strings.ToLower(faker.Adjective() + "-" + faker.Animal())
}

func checkFile(log *golog.Logger, url string) error {
	b, err := downloadGithubFile(url)
	if err != nil {
		return err
	}

	b2, err := readCache(url)
	// The case can be that the cached file does not exist yet.
	if err != nil {
		log.Info(err)
		b2 = []byte("")
	}

	if filesEqual(&b, &b2) {
		log.Info("Remote files match local cache")
		return nil
	}

	log.Info("Mismatch against local cache")
	nonce, err := gonanoid.New(8)
	if err != nil {
		return err
	}
	// Filename is just an MD5 hash of the manifest resource
	filename := toMD5Hash(url)
	tmpFilename := filename + "tmp"

	// Create a tmp file in case the deploy fails
	log.Info("Creating tmp compose file")
	if err = data.WriteManifestFile(tmpFilename, b); err != nil {
		return err
	}
	projectPrefix := getProjectNameFromHash(tmpFilename)
	newProjectName := projectPrefix + fmt.Sprintf("-%s", strings.ToLower(nonce))
	log.Info("Attempting redeploy")
	err = docker.StartComposeProject(tmpFilename, newProjectName)

	if err != nil {
		log.Error(err)
		return data.RemoveManifestFile(tmpFilename)
	}

	project, err := docker.GetOldProjectByPrefix(projectPrefix, newProjectName)
	if err != nil {
		log.Error(err)
		return data.RemoveManifestFile(tmpFilename)
	}

	if err = docker.RemoveComposeProject(filename, project.Name); err != nil {
		log.Error(err)
		return data.RemoveManifestFile(tmpFilename)
	}

	log.Info("Redeploy succeeded cleaning up...")
	if err = data.WriteManifestFile(filename, b); err != nil {
		return err
	}

	if err = data.RemoveManifestFile(tmpFilename); err != nil {
		return err
	}

	return nil
}

func checkFiles(log *golog.Logger) {
	urls, err := data.ReadManifestResources()

	if err != nil {
		log.Error(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		u := strings.TrimSpace(url)
		// Run every url request in parallell
		log.Info("Checking file ", url)
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
		log.Info("Polling disabled")
		return
	}

	log.Info("Polling enabled polling every", pollInterval, "minutes")
	interval, err := strconv.Atoi(pollInterval)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	for range ticker.C {
		checkFiles(log)
	}
}

func handleSync(ctx iris.Context) {
	checkFiles(ctx.Application().Logger())
	ctx.JSON(Msg{Message: "Successfully synced!"})
}

func RegisterSync(app *iris.Application) {
	syncApi := app.Party("/sync")
	syncApi.Post("/", handleSync)
	go startPoll(golog.Default)
}
