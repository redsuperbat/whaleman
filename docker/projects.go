package docker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/kataras/golog"
	"github.com/redsuperbat/whaleman/data"
)

func startAndPipeLogs(cmd *exec.Cmd, log *golog.Logger) error {
	cmdReader, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	startScanner(cmdReader, log.Info)
	return nil
}

func startScanner(reader io.ReadCloser, fn func(...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		m := scanner.Text()
		fn(m)
	}
}

func runCommand(log *golog.Logger, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if err := startAndPipeLogs(cmd, log); err != nil {
		return err
	}
	cmd.Wait()

	if cmd.ProcessState.ExitCode() == 1 {
		errMsg := fmt.Sprintf("unable to run command %v", cmd.Args)
		return errors.New(errMsg)
	}

	return nil
}

func StartComposeProject(filename string, project string) error {
	log := golog.New()
	log.SetPrefix(fmt.Sprintf("[%s] ", project))
	filepath := data.ManifestFilePath(filename)
	if err := runCommand(log, "docker-compose", "-f", filepath, "-p", project, "up", "-d"); err != nil {
		errMsg := fmt.Sprintf("unable to restart docker containers with manifest %s project %s", filepath, project)
		log.Error(errMsg)
		return err
	}

	return nil
}

func RemoveComposeProject(filename string, project string) error {
	log := golog.New()
	log.SetPrefix(fmt.Sprintf("[%s] ", project))
	filepath := data.ManifestFilePath(filename)

	if err := runCommand(log, "docker-compose", "-f", filepath, "-p", project, "down"); err != nil {
		errMsg := fmt.Sprintf("unable to remove containers with manifest %s project %s", filepath, project)
		log.Error(errMsg)
		return err
	}

	return nil
}

type Project struct {
	Name        string
	Status      string
	ConfigFiles string
}
