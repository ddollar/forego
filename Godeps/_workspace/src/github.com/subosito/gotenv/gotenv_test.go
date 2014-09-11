package gotenv

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var formats = []struct {
	in     string
	out    Env
	preset bool
}{
	// parses unquoted values
	{`FOO=bar`, Env{"FOO": "bar"}, false},

	// parses values with spaces around equal sign
	{`FOO =bar`, Env{"FOO": "bar"}, false},
	{`FOO= bar`, Env{"FOO": "bar"}, false},

	// parses double quoted values
	{`FOO="bar"`, Env{"FOO": "bar"}, false},

	// parses single quoted values
	{`FOO='bar'`, Env{"FOO": "bar"}, false},

	// parses escaped double quotes
	{`FOO="escaped\"bar"`, Env{"FOO": `escaped"bar`}, false},

	// parses empty values
	{`FOO=`, Env{"FOO": ""}, false},

	// expands variables found in values
	{"FOO=test\nBAR=$FOO", Env{"FOO": "test", "BAR": "test"}, false},

	// parses variables wrapped in brackets
	{"FOO=test\nBAR=${FOO}bar", Env{"FOO": "test", "BAR": "testbar"}, false},

	// reads variables from ENV when expanding if not found in local env
	{`BAR=$FOO`, Env{"BAR": "test"}, true},

	// expands undefined variables to an empty string
	{`BAR=$FOO`, Env{"BAR": ""}, false},

	// expands variables in quoted strings
	{"FOO=test\nBAR='quote $FOO'", Env{"FOO": "test", "BAR": "quote test"}, false},

	// does not expand escaped variables
	{`FOO="foo\$BAR"`, Env{"FOO": "foo$BAR"}, false},
	{`FOO="foo\${BAR}"`, Env{"FOO": "foo${BAR}"}, false},

	// parses yaml style options
	{"OPTION_A: 1", Env{"OPTION_A": "1"}, false},

	// parses export keyword
	{"export OPTION_A=2", Env{"OPTION_A": "2"}, false},

	// expands newlines in quoted strings
	{`FOO="bar\nbaz"`, Env{"FOO": "bar\nbaz"}, false},

	// parses varibales with "." in the name
	{`FOO.BAR=foobar`, Env{"FOO.BAR": "foobar"}, false},

	// strips unquoted values
	{`foo=bar `, Env{"foo": "bar"}, false}, // not 'bar '

	// ignores empty lines
	{"\n \t  \nfoo=bar\n \nfizz=buzz", Env{"foo": "bar", "fizz": "buzz"}, false},

	// ignores inline comments
	{"foo=bar # this is foo", Env{"foo": "bar"}, false},

	// allows # in quoted value
	{`foo="bar#baz" # comment`, Env{"foo": "bar#baz"}, false},

	// ignores comment lines
	{"\n\n\n # HERE GOES FOO \nfoo=bar", Env{"foo": "bar"}, false},

	// parses # in quoted values
	{`foo="ba#r"`, Env{"foo": "ba#r"}, false},
	{"foo='ba#r'", Env{"foo": "ba#r"}, false},

	// incorrect line format
	{"lol$wut", Env{}, false},
}

var fixtures = []struct {
	filename string
	results  Env
}{
	{
		"fixtures/exported.env",
		Env{
			"OPTION_A": "2",
			"OPTION_B": `\n`,
		},
	},
	{
		"fixtures/plain.env",
		Env{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "3",
			"OPTION_D": "4",
			"OPTION_E": "5",
		},
	},
	{
		"fixtures/quoted.env",
		Env{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": `\n`,
			"OPTION_E": "1",
			"OPTION_F": "2",
			"OPTION_G": "",
			"OPTION_H": "\n",
		},
	},
	{
		"fixtures/yaml.env",
		Env{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": `\n`,
		},
	},
}

func TestParse(t *testing.T) {
	for i, tt := range formats {
		if tt.preset {
			os.Setenv("FOO", "test")
		}

		exp := Parse(strings.NewReader(tt.in))

		x := fmt.Sprintf("%+v\n", exp)
		o := fmt.Sprintf("%+v\n", tt.out)

		if x != o {
			t.Logf("%q\n", tt.in)
			t.Errorf("(%d) %s != %s\n", i, x, o)
		}

		os.Clearenv()
	}
}

func TestLoad(t *testing.T) {
	for i, tt := range fixtures {
		Load(tt.filename)

		for key, val := range tt.results {
			if eval := os.Getenv(key); eval != val {
				t.Errorf("(%d) %s => %s != %s", i, key, eval, val)
			}
		}

		os.Clearenv()
	}
}

func TestLoadEnv(t *testing.T) {
	Load()

	tkey := "HELLO"
	val := "world"

	if tval := os.Getenv(tkey); tval != val {
		t.Errorf("%s => %s != %s", tkey, tval, val)
	}

	os.Clearenv()
}

func TestLoadNonExist(t *testing.T) {
	file := ".nonexist.env"

	err := Load(file)
	if err == nil {
		t.Errorf("Load(`%s`) => error: `no such file or directory` != nil", file)
	}
}
