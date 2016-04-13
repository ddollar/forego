## forego

<a href="https://circleci.com/gh/ddollar/forego">
  <img align="right" src="https://circleci.com/gh/ddollar/forego.svg?style=svg">
</a>

Foreman in Go.

### Installation

##### Downloads

Forego uses [Equinox](https://equinox.io) to automatically build binaries for various platforms.

See the [downloads page](https://dl.equinox.io/convox/forego/stable).

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
