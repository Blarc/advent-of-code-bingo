package utils

import (
	"log"
	"os"
)

func GetEnvVariable(name string) string {
	v, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("Variable %s not defined in .env file.\n", name)
	}

	return v
}
