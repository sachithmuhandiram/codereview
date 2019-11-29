FROM golang:alpine

RUN apk update

RUN apk add git

RUN set -e \
    && go get github.com/sirupsen/logrus \
    && go get github.com/google/uuid 

COPY . /go/src

WORKDIR /go/src

RUN go build -o main .

EXPOSE 7070

RUN chmod 755 main

CMD [ "./main" ]