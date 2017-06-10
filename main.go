package main

import (
	"log"

	"github.com/pborzenkov/vb/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
