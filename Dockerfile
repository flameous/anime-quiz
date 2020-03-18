FROM golang:1.14

WORKDIR /app

# Cache modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -v

EXPOSE 8080

CMD ["./anime-quiz"]