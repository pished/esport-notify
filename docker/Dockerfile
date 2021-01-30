# Build Environment
# ------------------------------
FROM golang:latest as build-env
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

Copy . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app.go


# Deployment Environment
# ------------------------------
FROM chromedp/headless-shell

RUN apt-get update && apt-get install ca-certificates -y

COPY --from=build-env /app/bin/app /app

ENTRYPOINT ["/app"]
