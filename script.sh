 #!/bin/bash

up() {
  docker-compose up
}

down() {
  docker-compose down
}

test() {
  go test -v -cover ./...
}

test2() {
  go test \
    ./common/...\
    ./controllers/...\
    ./middlewares/...\
    ./repositories/...\
    ./services/...\
    -coverprofile=coverage.txt -covermode=atomic
}

swag() {
  swag init -g app/app.go
}

run() {
  go run *.go
}