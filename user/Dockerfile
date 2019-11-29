FROM golang:alpine

RUN apk update
RUN apk add git


WORKDIR /go/src

RUN set -e \
    go get -u golang.org/x/crypto/... \ 
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/sirupsen/logrus

COPY . /go/src

RUN go build -o main .

EXPOSE 7071

RUN chmod 755 main

CMD [ "./main"]