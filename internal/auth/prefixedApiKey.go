package auth

import (
	apikey "github.com/joemiller/prefixed-api-key"
)

const prefix = "lshelf"

func NewPrefixedAPIKey() (*apikey.Key, error) {
	generator, err := apikey.NewGenerator(prefix)
	if err != nil {
		return nil, err
	}

	key, err := generator.GenerateAPIKey()

	return &key, err
}
