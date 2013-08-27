package main

import (
	"bufio"
	"fmt"
	"github.com/kr/pretty"
	"io"
	"os"
	"regexp"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type EnvEntry struct {
	name    string
	command string
}

type Env map[string]string

var _ = pretty.Println // lol

func ReadEnv(filename string) (Env, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return make(Env), nil
	}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return parseEnv(fd)
}

func parseEnv(r io.Reader) (Env, error) {
	env := make(Env)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := envEntryRegexp.FindStringSubmatch(scanner.Text())
		env[parts[1]] = parts[2]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Reading Env: %s", err)
	}
	return env, nil
}
