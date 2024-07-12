FROM golang:1.22.5-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o search-api .

EXPOSE 8080

CMD ["./search-api"]

