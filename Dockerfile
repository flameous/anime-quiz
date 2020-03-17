FROM golang:1.13.8

COPY . /go/src/github.com/flameous/anime-quiz
WORKDIR /go/src/github.com/flameous/anime-quiz

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["anime-quiz"]