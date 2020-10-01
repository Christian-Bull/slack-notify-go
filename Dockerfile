# gets latest golang image
FROM golang:alpine AS builder

# needed for go modules
ENV GO111MODULE=on
# ENV CGO_ENABLED=0 # not sure what this does yet

# # sets env vars for host
ARG TARGETOS
ARG TARGETARCH

# change working dir
WORKDIR /bin/app

# adds everything to the dir
ADD cmd /bin/app
ADD go.sum /bin/app
ADD go.mod /bin/app

# gets go modules
RUN go mod tidy -v
RUN go mod download
# builds app
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o main ./slack-notify/.

# run it yo
CMD ["/bin/app/main"]
