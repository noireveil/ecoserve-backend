FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o select-api cmd/api/main.go
EXPOSE 3000
CMD ["./select-api"]
