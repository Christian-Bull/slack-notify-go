# gets latest golang image
FROM golang:latest

# creates app dir
RUN mkdir /app

# adds everything to the dir
ADD . /app

# change working dir
WORKDIR /app

# builds app
RUN go build -o main ./cmd/slack-notify/.

# run it yo
CMD ["/app/main"]