# log-rush-distribution-server

A log distribution server implementing the log-rush datasource spec.

This service aims to be:
- an entry into the log-rush ecosystem
- an example implementation of the log-rush datsource spec
- the source code, on whose top the docs get generated

## Get Started

1. Install Go 1.18 (for instructions see [here](https://go.dev/doc/install))
2. Clone this repository 
   - ssh: `git clone git@github.com:log-rush/distribution-server.git` 
   - https: `git clone https://github.com/log-rush/distribution-server.git`
   - github cli: `gh repo clone log-rush/distribution-server`
3. Install dependencies: `go mod download`
4. Build the binary: `sh ./scripts/build.sh`
5. Run the binary `./log-rush-server`
(alt. using go run) `go run ./app/main.go`

A server will start on localhost:7000 with the swagger docs at [http://localhost:7000/swagger/index.html](http://localhost:7000/swagger/index.html)

### Generate Swagger Docs

Run `sh ./scripts/swagger.sh`
