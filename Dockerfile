FROM golang:1.21.0-alpine as golang_build

ENV APP_DIR /app
WORKDIR /app

VOLUME /app

RUN apk add build-base

RUN go env -w GOPATH=/app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download && go mod tidy

RUN go install github.com/githubnemo/CompileDaemon@latest

COPY ./ .

RUN go install github.com/onsi/ginkgo/v2/ginkgo
EXPOSE 3002

RUN ["chmod", "+x", "/app/entrypoint.sh"]
ENTRYPOINT ["./entrypoint.sh"]