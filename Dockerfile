# gets latest golang image
FROM golang:alpine

# install git
RUN apk update && apk add git

# creates app dir
RUN mkdir /app

# adds everything to the dir
ADD cmd /app

# change working dir
WORKDIR /app

# gets go pkgs - need to update to go modules
RUN go get -u github.com/slack-go/slack

# builds app
RUN go build -o main ./slack-notify/.

# run it yo
CMD ["/app/main"]