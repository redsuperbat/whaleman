package data

import (
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

func ensureFile(filepath string) {
	log := golog.New()
	if _, err := os.Stat(filepath); err == nil {
		return
	}
	if err := os.WriteFile(filepath, []byte(""), 0644); err != nil {
		log.Fatal(err)
	}
}

func appendToFile(filepath string, content string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(content + "\n"); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}

func ManifestFilePath(filename string) string {
	return toFilePath(filename)
}

func EnsureDataDir(log *golog.Logger) {
	log.Info("Ensuring ", dataDir(), " exists")
	if err := os.MkdirAll(dataDir(), 0700); err != nil {
		log.Fatal(err)
	}
	log.Info("Ensuring ", manifestResourceFile())
	ensureFile(manifestResourceFile())
}

func WriteManifestFile(filename string, content []byte) error {
	filepath := toFilePath(filename)
	return os.WriteFile(filepath, content, 0644)
}

func ManifestFileExists(filename string) bool {
	filepath := toFilePath(filename)
	if _, err := os.Stat(filepath); err != nil {
		return false
	}
	return true
}

func ReadManifestFile(filename string) ([]byte, error) {
	filepath := toFilePath(filename)
	return os.ReadFile(filepath)
}

func RemoveManifestFile(filename string) error {
	filepath := toFilePath(filename)
	return os.Remove(filepath)
}

func WriteManifestResource(url string) error {
	return appendToFile(manifestResourceFile(), url)
}

func ReadManifestResources() ([]string, error) {
	if b, err := os.ReadFile(manifestResourceFile()); err != nil {
		return nil, err
	} else {
		urls := slices.Filter(strings.Split(string(b), "\n"), func(s string) bool {
			return s != ""
		})
		return urls, nil
	}
}
