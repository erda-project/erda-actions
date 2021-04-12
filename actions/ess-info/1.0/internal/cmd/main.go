package main

import (
	"os"

	"github.com/erda-project/erda-actions/actions/ess-info/1.0/internal/ess"
)

func main() {
	err := ess.Process()
	if err != nil {
		os.Exit(1)
	}
}
