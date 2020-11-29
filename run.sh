#!/usr/bin/env bash

set -e -u

migrate() {
    if [[ $1 != '--down' ]]; then
        migrate -database "mysql://sa:qweqwe@tcp(localhost:3305)/hackernews_db" -path internal/pkg/db/migrations/mysql down
    else
        migrate -database "mysql://sa:qweqwe@tcp(localhost:3305)/hackernews_db" -path internal/pkg/db/migrations/mysql up
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