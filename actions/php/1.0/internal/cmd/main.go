package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/php/1.0/internal/pkg/build"
)

func main() {

	if err := build.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "PHP Action failed, err: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "PHP Action success\n")
	os.Exit(0)
}
