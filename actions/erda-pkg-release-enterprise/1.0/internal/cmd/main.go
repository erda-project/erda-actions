package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/erda-pkg-release-enterprise/1.0/internal/pkg"
)

func main() {
	if err := pkg.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
