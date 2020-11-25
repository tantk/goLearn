FROM golang:1.15.5-alpine3.12

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN go build ./main.go

CMD ["./main"]
# this is pointless because it cannot connect to mysql service
# need to setup user define network to make go middleware microservice connect to dockerized mysql
