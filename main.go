package main

import (
	"flag"
	"log"
	"os"

	"github.com/zercle/gofiber-skelton/config"
	server "github.com/zercle/gofiber-skelton/internal/infrastructure"
)

var (
	version string
	build   string
	runEnv  string
)

func main() {
	// Running flag
	if len(os.Getenv("ENV")) != 0 {
		runEnv = os.Getenv("ENV")
	} else {
		flagEnv := flag.String("env", "dev", "A config file name without .env")
		flag.Parse()
		runEnv = *flagEnv
	}
	if err := config.LoadConfig(runEnv); err != nil {
		log.Fatalf("error while loading the env:\n %+v", err)
	}

	server, err := server.NewServer(version, build, runEnv)
	if err != nil {
		log.Fatalf("error while create server:\n %+v", err)
	}

	server.Run()
}
