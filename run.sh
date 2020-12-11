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
    echo ${ENV_CONFIG} | base64 --decode > config.yml
    curl https://cli-assets.heroku.com/install-ubuntu.sh | sh
    heroku auth:token
    heroku container:login
    heroku container:push --recursive --app ${HEROKU_APP_NAME}
    heroku container:release web db --app ${HEROKU_APP_NAME}
    heroku logs --app ${HEROKU_APP_NAME}
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