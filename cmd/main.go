package main

import (
	"fmt"
	"os"

	"github.com/bashi/dotenv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s command [args...]\n", os.Args[0])
		os.Exit(1)
	}
	if err := dotenv.Run(os.Args[1], os.Args[2:]); err != nil {
		status := dotenv.ExitStatus(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(status)
	}
}
