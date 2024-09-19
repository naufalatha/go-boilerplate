# Pippin

GOLANG Boilerplate API

## Prerequisites:

Git clone base repository and install the following prerequisites:

1. Golang version 1.20 or above
2. Makefile latest version
3. Docker (optional)
4. Install jet `go install github.com/go-jet/jet/v2/cmd/jet@latest`
5. Install swaggo `go install github.com/swaggo/swag/cmd/swag@latest`

#### Development

_Note : Please ensure app.env file exist and configured properly before running the program._
First run `make jet-init` to generate models for communicating accross the code and then simply execte `make run` or run go compiler from source `go run cmd/main.go`. If you have docker installed, you can build the docker image using `docker compose up --build`. Please read the docker compose file for more details!

Git pull from branch main and checkout into a new branch, develop your feature inside the new branch. When developing a new feature, ensure test driven development is properly followed and write a clean and readable code. Document your feature using General API Annotation and generate swagger docs using `make swag-init`.

When creating a new merge request, be sure to provide description and environment variable if required. Keep merge request name relevant to the feature you are developing and delete the old feature branch when merged _(excluding the pipeline branch)_.

#### Command Glossary

`make run` run the service using app.env configuration
`make build` compile service using golang compiler into binary based on os kernel
`make build-docker` build docker using docker-compose file
`make swag-init` format and generate swagger documentation based on General API Annotation  
`make jet-init` generate code models based on url string connection
