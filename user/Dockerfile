FROM golang:alpine

RUN apk update
RUN apk add git

WORKDIR /go/src

RUN set -e \
    && go get -u golang.org/x/crypto/... \
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/sirupsen/logrus \
    && go get github.com/dgrijalva/jwt-go

# setting container time zone
# RUN set -e \
#   && apk add tzdata \
#   && cp /usr/share/zoneinfo/Asia/Colombo /etc/localtime \
#   && echo "Asia/Colombo" >  /etc/timezone \
#   && apk del tzdata

COPY . /go/src

RUN go build -o main .

EXPOSE 7071

RUN chmod 755 main

CMD [ "./main"]