## forego
[![Build Status](https://travis-ci.org/ddollar/forego.svg?branch=master)](https://travis-ci.org/ddollar/forego)

Foreman in Go.

### Installation

##### OS X (Homebrew)

    brew install forego

##### Precompiled Binaries

* Linux
  [[32-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-386/forego)]
  [[64-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-amd64/forego)]
* OS X
  [[32-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/darwin-386/forego)]
  [[64-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/darwin-amd64/forego)]
* Windows
  [[32-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/windows-386/forego.exe)]
  [[64-bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/windows-amd64/forego.exe)]

##### Compile from Source

    $ go get -u github.com/ddollar/forego

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO
