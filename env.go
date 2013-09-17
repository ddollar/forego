package main

import (
	"github.com/subosito/gotenv"
	"os"
	"regexp"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type Env map[string]string

func ReadEnv(filename string) (Env, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return make(Env), nil
	}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	env := make(Env)
	for key, val := range gotenv.Parse(fd) {
		env[key] = val
	}
	return env, nil
}
