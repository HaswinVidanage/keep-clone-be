#!/usr/bin/env bash

set -e

if [[ -z "${MYSQL_USER}" ]]; then
    source docker/dev/docker.env
fi


migrate() {
    connectionString="mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@localhost:3305/${MYSQL_DATABASE}"
    if [[ $1 != '--down' ]]; then
        migrate -database $connectionString -path internal/pkg/db/migrations/mysql down
    else
        migrate -database $connectionString -path internal/pkg/db/migrations/mysql up
    fi
}

genAll() {
    go generate ./...
}

test() {
    go test ./... -v
    echo "Test exited with exit code $?"
}

dockerBuild() {
 docker-compose build -t keep-app .
#   docker build  -f ./docker/prod/Go/Dockerfile -t keep-app .
}

dockerStop(){
    docker container stop $(docker container ls -aq)
    docker container rm $(docker container ls -aq)
}

if [[ $1 = 'migrate' ]]; then
    migrate $2
elif [[ $1 = 'gen-all' ]]; then
    genAll
elif [[ $1 = 'docker-build' ]]; then
    dockerBuild
elif [[ $1 = 'test' ]]; then
    test
else
    echo "Usage: ./run.sh (migrate|test|gen-all)"
    exit 2
fi