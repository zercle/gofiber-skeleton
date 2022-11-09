package main

import (
	"flag"
	"log"
	"os"

	"github.com/zercle/gofiber-skelton/configs"
	servers "github.com/zercle/gofiber-skelton/internal/infrastructure"
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
	if err := configs.LoadConfig(runEnv); err != nil {
		log.Panicf("error while loading the env:\n %+v", err)
	}

	server, err := servers.NewServer(version, build, runEnv)
	if err != nil {
		log.Panicf("error while create server:\n %+v", err)
	}

	server.Run()
}
