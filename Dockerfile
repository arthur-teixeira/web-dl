FROM golang:1.19.0

WORKDIR /go/src

COPY . .
RUN go mod download
