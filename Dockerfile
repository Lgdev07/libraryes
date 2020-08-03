FROM golang:alpine

WORKDIR /olist_challenge

COPY go.mod go.sum ./

RUN go mod download

COPY . .
