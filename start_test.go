package main

import (
	"os"
	"testing"
)

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

func TestPortFromEnv(t *testing.T) {
	env := make(Env)
	port, err := basePort(env)
	if err != nil {
		t.Fatalf("Can not get base port: %s", err)
	}
	if port != 5000 {
		t.Fatal("Base port should be 5000")
	}

	os.Setenv("PORT", "4000")
	port, err = basePort(env)
	if err != nil {
		t.Fatal("Can not get port: %s", err)
	}
	if port != 4000 {
		t.Fatal("Base port should be 4000")
	}

	env["PORT"] = "6000"
	port, err = basePort(env)
	if err != nil {
		t.Fatalf("Can not get base port: %s", err)
	}
	if port != 6000 {
		t.Fatal("Base port should be 6000")
	}

	env["PORT"] = "forego"
	port, err = basePort(env)
	if err == nil {
		t.Fatalf("Port 'forego' should fail: %s", err)
	}

}
