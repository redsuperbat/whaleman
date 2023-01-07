package data

import (
	"net/url"
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

func isValidUrl(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
}

func ManifestFilePath(filename string) string {
	return toFilePath(filename)
}

func EnsureDataDir(log *golog.Logger) {
	log.Info("Ensuring ", dataDir())
	if err := os.MkdirAll(dataDir(), 0700); err != nil {
		log.Fatal(err)
	}
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

func ReadManifestResources() []string {
	urls := strings.Split(os.Getenv("COMPOSE_FILE_RESOURCES"), ",")
	// remove invalid whitespace
	urls = slices.Map(urls, func(s string) string {
		return strings.TrimSpace(s)
	})
	return slices.Filter(urls, func(s string) bool {
		return isValidUrl(s)
	})
}
