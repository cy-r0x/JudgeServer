package main

import (
	"os"

	"github.com/judgenot0/judge-backend/cmd"
	"github.com/judgenot0/judge-backend/config"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}
	cmd.Serve(config.HttpPort)
}
