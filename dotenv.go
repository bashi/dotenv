package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	envFile                = ".env"
	defaultErrorExitStatus = 1
)

func parseLine(line string) (string, string, error) {
	fields := strings.SplitN(line, "=", 2)
	if len(fields) != 2 {
		return "", "", fmt.Errorf("Failed to parse line: %s", line)
	}
	return fields[0], fields[1], nil
}

func parseLines(r io.Reader) (string, string, error) {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", nil
	}
	if line[0] == '#' {
		return "", "", nil
	}
	return parseLine(line)
}

func setEnvFromReader(r io.Reader) error {
	for {
		key, value, err := parseLines(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// TODO: Don't expand when a value is enclosed by single quotes.
		value = os.ExpandEnv(value)
		os.Setenv(key, value)
	}
	return nil
}

func setEnvFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		// Nop when the given path doesn't exist.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	return setEnvFromReader(file)
}

func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run runs the given command with environment variables defined in .env file
// located in the current directory.
func Run(command string, args []string) error {
	err := setEnvFromFile(envFile)
	if err != nil {
		return err
	}
	return execute(command, args)
}

// ExitStatus returns an int value correspond to the given err, defaults to 1.
func ExitStatus(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return defaultErrorExitStatus
}
