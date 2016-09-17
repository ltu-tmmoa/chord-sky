# Chord Sky

LTU Project: Chord Sky (Mobile and Distributed Computing Systems)

THIS PROJECT IS UNDER ACTIVE DEVELOPMENT. THE DOCUMENTATION THAT COMES WITH IT
MAY REFLECT ITS INTENDED RATHER THAN ACTUAL FUNCTIONALITY.

As part of the D7024E course at LTU we produced _Chord Sky_, a distributed
key/value store. The store service is provided by one or more nodes that
together form a ring cluster. Values are distributed among the ring members
using the _Chord_ Distributed Hash Table (DHT) algorithm.

## Contributing

### Coding and Code Style

If the _Go_ language isn't known, the following resources may be used to learn
it.

- ["A Tour of Go"](https://tour.golang.org)
- ["How to Write Go Code"](https://golang.org/doc/code.html)
- ["Effective Go"](https://golang.org/doc/effective_go.html)

The last document may be regarded as a reference rather than a text to be read
as is. The _Go_ language comes with its own required code style guidelines,
which should be followed.

### Building and Running

Building and running is managed using the standard build tools that comes with
a regular [Go installation](https://golang.org/dl/).

### Development Environment

An editor that uses the official _Go_ static analysis and formatting tools
should be used.

#### The Atom Editor

If wanting to use the [Atom editor](https://atom.io), the `go-plus` plug-in is
the only one required to make the editor into a capable _Go_ IDE. It will
prompt you to download a set of related plug-ins once installed. Make sure that
_Go_ is installed and a suitable `GOPATH` is set before installing `go-plus`.

### Cloning the Repository

_Go_ uses a rather complicated way to manage code and its dependencies. In
order to make sure that no problems arise from using an unsuitable repository
folder location, please follow the following steps.

1. Make sure that a correct `GOPATH` is set using the following command:
```sh
$ go env
```
   It should produce a list containing the following keys, among many others.
   As long as the values are set things should be in order.
```sh
GOPATH="/home/user/projects/go"
GOROOT="/usr/local/go"
```
2. Download the repository using the following command:
```sh
$ go get github.com/ltu-tmmoa/chord-sky
```
3. Change the repository origin using the following commands:
```sh
$ cd $GOPATH/src/github.com/ltu-tmmoa/chord-sky
$ git remote rm origin
$ git remote add origin git@github.com:ltu-tmmoa/chord-sky.git
```
