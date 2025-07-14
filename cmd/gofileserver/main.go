package main

import (
	"os"

	_ "go.uber.org/automaxprocs"

	"github.com/pachirode/gofileserver/internal/gofileserver"
)

func main() {
	cmd := gofileserver.NewGofileserverCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
