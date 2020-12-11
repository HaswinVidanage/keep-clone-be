FROM golang:1.13.7

# Force the go compiler to use modules
ENV GO111MODULE=on

WORKDIR ./app
RUN mkdir build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN ls

# Build the Go app
RUN go build -o ./build/keep_be .

RUN ls /go/bin/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./build/keep_be"]