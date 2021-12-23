package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/email/1.0/internal/pkg/build"
)

func main() {
	err := build.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
