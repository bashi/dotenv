package main

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
	env_filename           = ".env"
	defaultErrorExitStatus = 1
)

func parseLine(r *bufio.Reader) (string, string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	line = strings.TrimSpace(line)
	if line[0] == '#' {
		return "", "", nil
	}
	fields := strings.Split(line, "=")
	if len(fields) != 2 {
		return "", "", fmt.Errorf("Failed to parse line: %s", line)
	}
	return fields[0], fields[1], nil
}

func setEnvFromReader(r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		key, value, err := parseLine(br)
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

func execute(name string, args []string) error {
	cmd := exec.Command(name, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		_, err := io.Copy(stdin, os.Stdin)
		if err != nil {
			panic(err)
		}
	}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run(name string, args []string) error {
	err := setEnvFromFile(env_filename)
	if err != nil {
		return err
	}
	return execute(name, args)
}

func exitStatus(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return defaultErrorExitStatus
}

func main() {
	if err := run(os.Args[1], os.Args[2:]); err != nil {
		status := exitStatus(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(status)
	}
}
