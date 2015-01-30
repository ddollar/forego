## forego

<div style="float:right">
  <a href="https://travis-ci.org/ddollar/forego">
    <img src="https://travis-ci.org/ddollar/forego.svg?branch=master">
  </a>
</div>

Foreman in Go.

### Installation

##### OS X (Homebrew)

    brew install forego

##### Precompiled Binaries

* [Linux](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-amd64/forego)
* [OSX](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/darwin-amd64/forego)
* [Windows](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/windows-amd64/forego.exe)

##### Compile from Source

    $ go get -u github.com/ddollar/forego

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO

### License

Apache 2.0 &copy; 2015 David Dollar
