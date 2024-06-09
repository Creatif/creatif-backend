FROM golang:1.22.3-alpine as golang_build

ENV APP_DIR /app
WORKDIR /app

RUN apk add build-base

COPY go.mod .
COPY go.sum .

RUN go install github.com/air-verse/air@latest
RUN go get github.com/onsi/ginkgo/v2/ginkgo
RUN go install github.com/onsi/ginkgo/v2/ginkgo

RUN go mod download
RUN go mod tidy

COPY . .

EXPOSE 3002

CMD ["air", "-c", "/app/cmd/http/.air.toml"]