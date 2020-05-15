package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func addBasicAuth(req *http.Request, usernameEnv string, passwordEnv string) (err error) {
	username, ok := os.LookupEnv(usernameEnv)
	if !ok {
		err = missingEnv(usernameEnv)
		return
	}

	password, ok := os.LookupEnv(passwordEnv)
	if !ok {
		err = missingEnv(passwordEnv)
		return
	}

	req.SetBasicAuth(username, password)
	return
}

func missingEnv(env string) error {
	return errors.New(fmt.Sprintf("%s not specified!", env))
}
