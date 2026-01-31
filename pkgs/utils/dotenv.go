package utils

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var FlagSet *AppFlagSet

func ParseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Loading env(%s): %s", envFile, err.Error())
		}
	} else if envFile != ".env" {
		log.Fatalf("Loading env(%s): %s", envFile, err.Error())
	}

	FlagSet = NewFlagSet("app", flag.CommandLine)

	FlagSet.Parse([]string{})
}
