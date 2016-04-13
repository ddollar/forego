# gotenv

Load environment variables dynamically in Go.

|-              | -                                                  |
|---------------|----------------------------------------------------|
| Build Status  | [![Build Status][drone-img]][drone-url]            |
| Coverage      | [![Coverage Status][coveralls-img]][coveralls-url] |
| Documentation | http://godoc.org/github.com/subosito/gotenv        |

## Installation

```bash
$ go get github.com/subosito/gotenv
```

## Usage

Store your configuration to `.env` file on your root directory of your project:

```
APP_ID=1234567
APP_SECRET=abcdef
```

Put the gotenv package on your `import` statement:

```go
import "github.com/subosito/gotenv"
```

Then somewhere on your application code, put:

```go
gotenv.Load()
```

Behind the scene it will then load `.env` file and export the valid variables to the environment variables. Make sure you call the method as soon as possible to ensure all variables are loaded, say, put it on `init()` function.

Once loaded you can use `os.Getenv()` to get the value of the variable.

Here's the final example:

```go
package main

import (
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	gotenv.Load()
}

func main() {
	log.Println(os.Getenv("APP_ID"))     // "1234567"
	log.Println(os.Getenv("APP_SECRET")) // "abcdef"
}
```

You can also load other than `.env` file if you wish. Just supply filenames when calling `Load()`:

```go
gotenv.Load(".env.production", "credentials")
```

That's it :)

### Another Scenario

Just in case you want to parse environment variables from any `io.Reader`, gotenv keeps its `Parse()` function as public API so you can utilize that.

```go
// import "strings"

pairs := gotenv.Parse(strings.NewReader("FOO=test\nBAR=$FOO"))
// gotenv.Env{"FOO": "test", "BAR": "test"}

pairs = gotenv.Parse(strings.NewReader(`FOO="bar"`))
// gotenv.Env{"FOO": "bar"}
```

Parse ignores invalid lines and returns `Env` of valid environment variables.

### Formats

The gotenv supports various format for defining environment variables. You can see more about it on:

- [fixtures](./fixtures)
- [gotenv_test.go](./gotenv_test.go)

## Notes

The gotenv package is a Go port of [`dotenv`](https://github.com/bkeepers/dotenv) project. Most logic and regexp pattern is taken from there and aims will be compatible as close as possible.

[drone-img]: https://drone.io/github.com/subosito/gotenv/status.png
[drone-url]: https://drone.io/github.com/subosito/gotenv/latest
[coveralls-img]: https://coveralls.io/repos/subosito/gotenv/badge.png?branch=master
[coveralls-url]: https://coveralls.io/r/subosito/gotenv?branch=master

