FROM golang:alpine

WORKDIR /libraryes

COPY go.mod go.sum ./

RUN go mod download

COPY . .
