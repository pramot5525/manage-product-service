FROM golang:1.26

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY docs ./docs
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
