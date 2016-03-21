package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/subosito/gotenv"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type Env map[string]string

type envFiles []string

func (e *envFiles) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *envFiles) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func loadEnvs(files []string) (Env, error) {
	if len(files) == 0 {
		files = []string{".env"}
	}

	env := make(Env)

	for _, file := range files {
		tmpEnv, err := ReadEnv(file)

		if err != nil {
			return nil, err
		}

		// Merge the file I just read into the env.
		for k, v := range tmpEnv {
			env[k] = v
		}
	}
	return env, nil
}

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

func (e *Env) asArray() (env []string) {
	for _, pair := range os.Environ() {
		env = append(env, pair)
	}
	for name, val := range *e {
		env = append(env, fmt.Sprintf("%s=%s", name, val))
	}
	return
}
