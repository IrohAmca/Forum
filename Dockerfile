FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apk add --no-cache build-base

COPY . .

ENV CGO_ENABLED=1 
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
