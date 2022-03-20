package main

import (
	_ "embed"
	"os"

	"github.com/julis-bw-nw/julis-service-register-app/cmd"
)

//go:embed LICENSE
var license string

func main() {
	if err := cmd.Execute(license); err != nil {
		os.Exit(1)
	}
}
