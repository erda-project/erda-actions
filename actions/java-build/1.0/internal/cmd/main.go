package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/java-build/1.0/internal/pkg/build"
)

func main() {
	err := build.Execute()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
