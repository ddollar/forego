package main

import "testing"

func TestParseConcurrencyFlagEmpty(t *testing.T) {
	c, err := parseConcurrency("")
	if err != nil {
		t.Fatal(err)
	}
	if len(c) > 0 {
		t.Fatal("expected no concurrency settings with ''")
	}
}

func TestParseConcurrencyFlagSimle(t *testing.T) {
	c, err := parseConcurrency("foo=2")
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 1 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}
}

func TestParseConcurrencyFlagMultiple(t *testing.T) {
	c, err := parseConcurrency("foo=2,bar=3")
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 2 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}

	if c["bar"] != 3 {
		t.Fatal("expected concurrency settings of 3 with 'bar=3'")
	}
}

func TestParseConcurrencyFlagNonInt(t *testing.T) {
	_, err := parseConcurrency("foo=x")
	if err == nil {
		t.Fatal("foo=x should fail")
	}
}

func TestParseConcurrencyFlagWhitespace(t *testing.T) {
	c, err := parseConcurrency("foo   =   2, bar = 3")
	if err != nil {
		t.Fatalf("foo   =   2, bar = 4 should not fail:%s", err)
	}

	if len(c) != 2 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}

	if c["bar"] != 3 {
		t.Fatal("expected concurrency settings of 3 with 'bar=3'")
	}
}

func TestParseConcurrencyFlagMultipleEquals(t *testing.T) {
	_, err := parseConcurrency("foo===2")
	if err == nil {
		t.Fatalf("foo===2 should fail: %s", err)
	}
}

func TestParseConcurrencyFlagNoValue(t *testing.T) {
	_, err := parseConcurrency("foo=")
	if err == nil {
		t.Fatalf("foo= should fail: %s", err)
	}

	_, err = parseConcurrency("=")
	if err == nil {
		t.Fatalf("= should fail: %s", err)
	}

	_, err = parseConcurrency("=1")
	if err == nil {
		t.Fatalf("= should fail: %s", err)
	}

	_, err = parseConcurrency(",")
	if err == nil {
		t.Fatalf(", should fail: %s", err)
	}

	_, err = parseConcurrency(",,,")
	if err == nil {
		t.Fatalf(",,, should fail: %s", err)
	}

}

func TestMultipleEnvironmentFiles(t *testing.T) {
	envs := []string{"fixtures/envs/.env1", "fixtures/envs/.env2"}
	env, err := parseEnvironment(envs)

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
