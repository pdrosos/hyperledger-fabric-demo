# Customer App

## Installation

Clone the application: `git clone https://github.com/pdrosos/hyperledger-fabric-demo && cd hyperledger-fabric-demo/customer`

Start the docker containers: `docker-compose up -d`

Application runs on http://localhost:7777


## For developers

### Download
To download the application for development as a Go package inside your `$GOPATH`, run `go get github.com/pdrosos/hyperledger-fabric-demo/customer`.

### Local development
To develop and debug the application locally you need [Go 1.10](https://golang.org/) and [Dep](https://golang.github.io/dep/) dependency manager installed locally on your machine.

### Install dependencies
Application uses dep to store its dependencies inside the `vendor` directory. <br>
To install the application dependencies, run or `dep ensure`

### Format code
Run `make fmt` to format your code according to the Go style guide. Do this before every commit.

### Install
Rebuild the application from source: `make`

### Tests
Run the tests: `make test`

### Docker container to run the application
Application uses Golang 1.10. It can run inside the `app` container. <br>
To start the containers, go to the application directory: `cd $GOPATH/src/github.com/pdrosos/hyperledger-fabric-demo/customer` and run `docker-compose up -d`