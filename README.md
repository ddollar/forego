## forego
[![Build Status](https://travis-ci.org/ddollar/forego.svg?branch=master)](https://travis-ci.org/ddollar/forego)

Foreman in Go.

### Installation

##### OSX Homebrew

    brew install forego

##### Precompiled Binaries

Download from [gobuild.io](http://gobuild.io/github.com/ddollar/forego)

##### Compile from Source

    $ go get -u github.com/ddollar/forego

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO
