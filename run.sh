#!/usr/bin/env bash

set -e

if [[ -z "${MYSQL_USER}" ]]; then
    source ./docker/docker.env
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

if [[ $1 = 'migrate' ]]; then
    migrate $2
elif [[ $1 = 'gen-all' ]]; then
    genAll
elif [[ $1 = 'test' ]]; then
    test
else
    echo "Usage: ./run.sh (migrate|test|gen-all)"
    exit 2
fi