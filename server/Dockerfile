FROM golang:1.14.2-alpine3.11

WORKDIR /go/src/github.com/benrhyshoward/hitthespot/server

COPY . .

RUN go install -v

EXPOSE 8080

CMD ["hitthespot"]