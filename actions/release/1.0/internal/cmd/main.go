package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/pkg"
)

func main() {
	if err := pkg.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println(fmt.Sprintf("err info: %v", err))
		os.Exit(1)
	}
}
