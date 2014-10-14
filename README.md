## forego
[![Build Status](https://travis-ci.org/ddollar/forego.svg?branch=master)](https://travis-ci.org/ddollar/forego)

Foreman in Go.

### Installation

##### OSX Homebrew

    brew install forego

##### Precompiled Binaries

* [Linux 32bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-386/forego)
* [Linux 64bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-amd64/forego)
* [OS X 32bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/darwin-386/forego)
* [OS X 64bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/darwin-amd64/forego)
* [Windows 32bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/windows-386/forego.exe)
* [Windows 64bit](https://godist.herokuapp.com/projects/ddollar/forego/releases/current/windows-amd64/forego.exe)

##### Compile from Source

    $ go get -u github.com/ddollar/forego

### Usage

Forego is compatible with the Procfile file format.  This is also what the foreman Ruby gem uses.  There's documentation [here](http://ddollar.github.io/foreman/#PROCFILE) but here's a simple shell example that should get you started.  Create a file named Procfile in your project's top-level directory.

    some_server: while [ true ]; do sleep 10; done
    another_process: while [ true ]; do sleep 10; done

When you are ready to start your processes defined in Procfile run:

    $ forego start

Forego also comes with command line help about starting processes:

    $ forego help start

For more help, see `forego help`

