FROM golang:1.25.3-alpine

WORKDIR /app

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o main example/main.go

CMD ["./main"]