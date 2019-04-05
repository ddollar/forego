## Forego

[Foreman](https://github.com/ddollar/foreman) in Go.

## Installing

For your convenience, Forego can be installed via Homebrew. However you cannot install it along side the original Forego, from which this was forked.

Install the latest version thusly:

```
$ brew install 'http://bww-artifacts.s3.amazonaws.com/forego/latest/forego.rb'
```

Upgrade (or uninstall) it as follows:

```
$ brew upgrade 'http://bww-artifacts.s3.amazonaws.com/forego/latest/forego.rb'
```

Uninstall as follows:

```
$ brew uninstall forego
```

### Usage

    $ cat Procfile
    web: bin/web start -p $PORT
    worker: bin/worker queue=FOO

    $ forego start
    web    | listening on port 5000
    worker | listening to queue FOO

Use `forego help` to get a list of available commands, and `forego help
<command>` for more detailed help on a specific command.

### License

Apache 2.0 &copy; 2015 David Dollar, 2018 Brian W. Wolter
