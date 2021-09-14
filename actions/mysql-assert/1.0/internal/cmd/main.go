package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/mysql-assert/1.0/internal/pkg/build"
)

func main() {
	if err := build.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
