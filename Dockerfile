# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY internal/database/* ./internal/database/
COPY media.db ./
RUN ls -altr ./
RUN go build -o /docker-go-socialmedia
RUN ls -altr /
EXPOSE 8080
CMD [ "/docker-go-socialmedia"]