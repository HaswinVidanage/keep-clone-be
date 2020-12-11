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
    base64 --decode -i ${ENV_CONFIG} -o config.yml
    curl https://cli-assets.heroku.com/install-ubuntu.sh | sh
    heroku auth:token
    heroku container:push --recursive
    heroku container:release web db
    heroku logs
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