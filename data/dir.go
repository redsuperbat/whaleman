package data

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kataras/golog"
	"github.com/redsuperbat/whaleman/slices"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func dataDir() string {
	return getEnv("DATA_DIR", "/var/lib/whaleman")
}

func toFilePath(filename string) string {
	return dataDir() + "/" + filename
}

func manifestResourceFile() string {
	return dataDir() + "/resources"
}

func EnsureDataDir(log *golog.Logger) {
	log.Debug("Ensuring ", dataDir(), " exists")
	if err := os.MkdirAll(dataDir(), 0700); err != nil {
		log.Fatal(err)
	}

	path := manifestResourceFile()
	log.Debug("Ensuring ", path, " exists")
	if _, err := os.Stat(path); err == nil {
		return
	}
	log.Debug(path, " does not exist initializing an empty one")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
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

func RemoveManifestFile(filename string) error {
	filepath := toFilePath(filename)
	if err := os.Remove(filepath); err != nil {
		return err
	}
	return nil
}

func WriteManifestResource(url string) error {
	f, err := os.OpenFile(manifestResourceFile(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(url + "\n"); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}

func ReadManifestResources() (error, []string) {
	if b, err := ioutil.ReadFile(manifestResourceFile()); err != nil {
		return err, nil
	} else {
		urls := slices.Filter(strings.Split(string(b), "\n"), func(s string) bool {
			return s != ""
		})
		return nil, urls
	}
}
