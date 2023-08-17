FROM golang:1.21.0-alpine as golang_build

ENV APP_DIR /app
WORKDIR /app

VOLUME /app

RUN apk add build-base
RUN go env -w GOPATH=/app

COPY . .

RUN go mod download && go mod tidy

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/onsi/ginkgo/v2/ginkgo

RUN cd /app/cmd/http && go build -o go_api_build

#CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/cmd/http/go_api_build"]

RUN ["chmod", "+x", "/app/entrypoint.sh"]
ENTRYPOINT ["./entrypoint.sh"]