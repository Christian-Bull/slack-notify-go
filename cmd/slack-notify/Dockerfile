# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Christian Bull <bullcs95@gmail.com>"

# creates app dir
RUN mkdir /app
ADD . /app/

# Set the Current Working Directory inside the container
WORKDIR /app

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main