## Up

Monitor your web applications. RESTful administration interface.

### Configuration

    $ export PORT=4242
    $ export SECRET=foobar

### Build

    $ go build

### Running

    $ go build
    $ ./up-server

Then navigate to http://localhost:4242 in your web browser.

### Test

    $ go test

### Build Dependencies

You'll need [gccgo](http://golang.org/doc/install).

To install all of the necessary dependencies:

    $ cd up-server/
    $ go get .
    ...
