package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ddollar/forego/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/ddollar/forego/Godeps/_workspace/src/github.com/subosito/gotenv"
)

var envEntryRegexp = regexp.MustCompile("^([A-Za-z_0-9]+)=(.*)$")

type Env map[string]string

type envFiles []string

func (e *envFiles) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *envFiles) Set(value string) error {
	*e = append(*e, fullPath(value))
	return nil
}

func fullPath(file string) string {
	root := filepath.Dir(".")
	return filepath.Join(root, file)
}

func loadEnvs(files []string, etcdHost string, etcdKey string) (Env, error) {
	var err error
	env := make(Env)

	if len(files) == 0 {
		env, err = ReadEnv(fullPath(".env"))
		if err != nil {
			return nil, err
		}
	} else {
		// Handle multiple environment files
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
	}

	// Get etcd configuration from env if not yet defined
	if etcdHost == "" {
		etcdHost = env["ETCD_HOST"]
	}
	if etcdKey == "" {
		etcdKey = env["ETCD_KEY"]
	}
	// Read from etcd
	if etcdHost != "" {
		etcdEnv, err := ReadEnvFromEtcd(etcdHost, etcdKey)
		if err != nil {
			return nil, err
		}

		// Merge etcd env
		for k, v := range etcdEnv {
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

func ReadEnvFromEtcd(host, key string) (Env, error) {
	etcdClient := etcd.NewClient([]string{host})

	result, err := etcdClient.Get(key, true, true)
	if err != nil {
		return nil, err
	}

	env := make(Env)
	for _, node := range result.Node.Nodes {
		k := strings.Split(node.Key, "/")
		env[strings.ToUpper(k[len(k)-1])] = node.Value
	}
	return env, nil
}
