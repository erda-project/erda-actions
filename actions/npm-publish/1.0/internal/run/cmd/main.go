package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/npm-publish/1.0/internal/npm"
	"github.com/erda-project/erda-actions/actions/npm-publish/1.0/internal/run"
)

func main() {
	NPM := npm.NewNPM()
	command := run.NewCommand(NPM)

	var request run.Request
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal("reading request from stdin", err)
	}

	err := checkParams(request)
	if err != nil {
		fatal("parameter required", err)
	}

	//request.Params.Path = filepath.Join(os.Args[1], request.Params.Path)
	response, err := command.Run(request)
	if err != nil {
		fatal("running command", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fatal("writing response to stdout", err)
	}
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func checkParams(request run.Request) error {
	var err error
	if request.Params.UserName == "" {
		err = errors.New("username")
	}
	if request.Params.Password == "" {
		err = errors.New("password")
	}
	if request.Params.Email == "" {
		err = errors.New("email")
	}
	if request.Params.Path == "" {
		err = errors.New("path")
	}
	return err
}
