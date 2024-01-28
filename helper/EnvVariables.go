package helpers

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(variables []string) ([]string, error) {

	if len(variables) == 0 {
		return nil, errors.New("No variables passed to GetEnv Function")
	}

	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("Trouble loading env file")
	}

	var envVariables []string

	for _, variable := range variables {
		envVariables = append(envVariables, os.Getenv(variable))
	}
	return envVariables, nil
}
