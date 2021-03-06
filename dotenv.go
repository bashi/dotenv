package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

const (
	envFile                = ".env"
	defaultErrorExitStatus = 1
)

func parseLine(line string, prefix string) (string, string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", nil
	}
	if line[0] == '#' {
		return "", "", nil
	}
	fields := strings.SplitN(line, "=", 2)
	if len(fields) != 2 {
		return "", "", fmt.Errorf("Failed to parse line: %s", line)
	}

	key := fields[0]
	value := fields[1]
	// Quick hack: Treat |value| as a relative path when it starts with
	// "./" or "../"
	if strings.HasPrefix(value, "./") || strings.HasPrefix(value, "../") {
		value = path.Clean(prefix + value)
	}
	// TODO: Don't expand when a value is enclosed by single quotes.
	value = os.ExpandEnv(value)
	return key, value, nil
}

func setEnvFromReader(r io.Reader, prefix string) error {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		key, value, err := parseLine(line, prefix)
		if err != nil {
			return err
		}
		os.Setenv(key, value)
	}
	return nil
}

func setEnvFromFile(path string, prefix string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return setEnvFromReader(file, prefix)
}

func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findEnvFilePath() (string, string, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	prefix := ""
	for {
		cur := path.Join(dirname, envFile)
		if _, err := os.Stat(cur); err == nil {
			return cur, prefix, nil
		}
		i := strings.LastIndex(dirname, "/")
		if i < 0 {
			break
		}
		dirname = dirname[:i]
		prefix = prefix + "../"
	}
	return "", "", fmt.Errorf("Cannot find %s file", envFile)
}

// Run runs the given command with environment variables defined in .env file
// located in the current directory.
func Run(command string, args []string) error {
	pathname, prefix, err := findEnvFilePath()
	if err != nil {
		return err
	}
	err = setEnvFromFile(pathname, prefix)
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
