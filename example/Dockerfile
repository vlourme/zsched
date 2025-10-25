FROM golang:1.25.3-alpine

WORKDIR /app

COPY . .

RUN go build -o main example/main.go

CMD ["./main"]