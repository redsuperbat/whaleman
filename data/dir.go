package data

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func manifestVersionFile() string {
	return dataDir() + "/versions"
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

func IncrementManifestVersion(filename string) error {
	currentVersion := GetManifestVersion(filename)
	filepath := manifestVersionFile()
	if currentVersion == -1 {
		return errors.New("file does not exist")
	}
	input, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, filename) {
			lines[i] = filename + fmt.Sprintf("%v", currentVersion+1)
		}
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetManifestVersion(filename string) int {
	filepath := manifestVersionFile()
	reader, err := os.Open(filepath)
	if err != nil {
		return -1
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, filename) {
			continue
		}
		if i, err := strconv.Atoi(line[len(filename):]); err != nil {
			return -1
		} else {
			return i
		}
	}
	return -1
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
	log.Info("Ensuring ", manifestVersionFile())
	ensureFile(manifestVersionFile())
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
