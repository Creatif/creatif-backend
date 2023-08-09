#!/bin/sh
_term() {
  kill -TERM "$child" 2>/dev/null
}

trap _term SIGTERM

/go/bin/CompileDaemon --build="go build -o go_api_build" --command=/app/cmd/http/./go_api_build --directory="/app/cmd/http" --pattern="(.+\.go)$" --polling --polling-interval=1000 --graceful-kill=true --color=true &
child=$!
wait "$child"