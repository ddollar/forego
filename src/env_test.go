package main

import "testing"

func TestMultipleEnvironmentFiles(t *testing.T) {
	envs := []string{"../fixtures/envs/.env1", "../fixtures/envs/.env2"}
	env, err := loadEnvs(envs)

	if err != nil {
		t.Fatalf("Could not read environments: %s", err)
	}

	if env["env1"] == "" {
		t.Fatalf("$env1 should be present and is not")
	}

	if env["env2"] == "" {
		t.Fatalf("$env2 should be present and is not")
	}
}
