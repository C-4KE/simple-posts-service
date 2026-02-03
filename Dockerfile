FROM golang:1.25.6

WORKDIR /app

copy . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]