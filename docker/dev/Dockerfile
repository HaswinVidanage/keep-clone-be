FROM golang:1.13.7

# Force the go compiler to use modules
ENV GO111MODULE=on

RUN mkdir /usr/src/keep_be
WORKDIR /usr/src/keep_be

# We want to populate the module cache based on the go.{mod,sum} files.
#COPY ../go.mod .
#COPY ../go.sum .
#
##This is the ‘magic’ step that will download all the dependencies that are specified in
## the go.mod and go.sum file.
## Because of how the layer caching system works in Docker, the go mod download
## command will _ only_ be re-run when the go.mod or go.sum file change
## (or when we add another docker instruction this line)
#RUN go mod download
#
## Here we copy the rest of the source code
#COPY .. .
#
## And compile the project
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
#
#RUN apt-get update
#RUN apt-get install vim -y
#RUN go get "github.com/go-sql-driver/mysql"
#
#
#EXPOSE 8080