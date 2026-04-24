FROM golang:1.25-alpine

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ecoserve-api cmd/api/main.go
CMD ["./ecoserve-api"]