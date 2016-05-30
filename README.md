## forego

<a href="https://circleci.com/gh/ddollar/forego">
  <img align="right" src="https://circleci.com/gh/ddollar/forego.svg?style=svg">
</a>

Foreman in Go.

### Installation

[Downloads](https://dl.equinox.io/ddollar/forego/stable)

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
