#!/usr/bin/env bash

set -e

migrate() {
    if [[ -z "${MYSQL_USER}" ]]; then
        source docker.env
    fi
    connectionString="mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@localhost:3305/${MYSQL_DATABASE}"
    if [[ $1 != '--down' ]]; then
        migrate -database $connectionString -path internal/pkg/db/migrations/mysql down
    else
        migrate -database $connectionString -path internal/pkg/db/migrations/mysql up
    fi
}

deploy() {
    echo ${ENV_CONFIG}
    echo ${ENV_CONFIG} | base64 --decode > config.yml
    curl https://cli-assets.heroku.com/install-ubuntu.sh | sh
    HEROKU_API_KEY=${HEROKU_API_KEY} heroku auth:token
    HEROKU_API_KEY=${HEROKU_API_KEY} heroku container:push --recursive
    HEROKU_API_KEY=${HEROKU_API_KEY} heroku container:release web db
    HEROKU_API_KEY=${HEROKU_API_KEY} heroku logs
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
elif [[ $1 = 'deploy' ]]; then
    deploy
elif [[ $1 = 'test' ]]; then
    test
else
    echo "Usage: ./run.sh (migrate|test|gen-all|deploy)"
    exit 2
fi