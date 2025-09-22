package main

import (
	"log"
	"os"

	env "github.com/joho/godotenv"
	"github.com/judgenot0/judge-backend/cmd"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatalln("ENV Not Found")
		return
	}
	HTTP_PORT := os.Getenv("HTTP_PORT")
	cmd.Serve(HTTP_PORT)
}
