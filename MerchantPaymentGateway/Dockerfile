



# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Oluwasey Otun <seyijava@gmail.com>"

# Set the Current Working Directory inside the container

RUN mkdir /app

# We copy everything in the root directory
# into our /app directory
ADD . /app


WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download



# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 9090

# Command to run the executable
CMD ["./main"]