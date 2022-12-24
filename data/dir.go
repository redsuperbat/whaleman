package data

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kataras/golog"
)

const (
	DATA_DIR               = "/var/lib/whaleman"
	MANIFEST_RESOURCE_FILE = DATA_DIR + "/resources"
)

func toFilePath(filename string) string {
	return DATA_DIR + "/" + filename
}

func EnsureDataDir(log *golog.Logger) {
	log.Debug("Ensuring", DATA_DIR, "exists")
	if err := os.MkdirAll(DATA_DIR, 0700); err != nil {
		log.Fatal(err)
	}
	log.Debug("Ensuring", MANIFEST_RESOURCE_FILE, "exists")
	if _, err := os.Stat(MANIFEST_RESOURCE_FILE); err == nil {
		return
	}

	log.Debug("File does not exist initializing an empty one")
	if err := os.WriteFile(MANIFEST_RESOURCE_FILE, []byte(""), 0644); err != nil {
		log.Fatal(err)
	}
}

func WriteManifestFile(filename string, content []byte) error {
	filepath := toFilePath(filename)
	if err := ioutil.WriteFile(filepath, content, 0644); err != nil {
		return err
	}
	return nil
}

func ManifestFileExists(filename string) bool {
	filepath := toFilePath(filename)
	if _, err := os.Stat(filepath); err != nil {
		return false
	}
	return true
}

func ReadManifestFile(filename string) (error, []byte) {
	filepath := toFilePath(filename)
	if b, err := ioutil.ReadFile(filepath); err != nil {
		return err, nil
	} else {
		return nil, b
	}
}

func WriteManifestResource(url string) error {
	f, err := os.OpenFile("text.log", os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(url + "\n"); err != nil {
		return err
	}
	return nil
}

func ReadManifestResources() (error, []string) {
	if b, err := ioutil.ReadFile(MANIFEST_RESOURCE_FILE); err != nil {
		return err, nil
	} else {
		return nil, strings.Split(string(b), "\n")
	}
}
