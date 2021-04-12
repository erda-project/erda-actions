package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/lib-publish/1.0/internal/pkg/build"
)

func main() {
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
