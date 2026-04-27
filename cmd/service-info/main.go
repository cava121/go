package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err);
		os.Exit(1);
	}
}

func run() error {
	cmdName := os.Args[1]
	os.Args = os.Args[1:]

	switch cmdName {
		case "versions":
			return runVersions();
		case "is-deployd":
			// TODO
	}

	return  nil
}

func runVersions() error {
	res, err := http.Get("http://localhost:8090/debug/info")
	if err != nil {
		return fmt.Errorf("get debug info %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body);
	if err != nil {
		return fmt.Errorf("read data %w", err)

	}

	fmt.Println(strings.TrimSpace(string(data)));
	return  nil;
}
