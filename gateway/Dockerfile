FROM golang:alpine

RUN apk update

RUN apk add git

RUN set -e \
    && go get github.com/sirupsen/logrus \
    && go get github.com/google/uuid \
    && go get github.com/go-sql-driver/mysql

# # setting container time zone
# RUN set -e \
#   && apk add tzdata \
#   && cp /usr/share/zoneinfo/Asia/Colombo /etc/localtime \
#   && echo "Asia/Colombo" >  /etc/timezone \
#   && apk del tzdata

COPY . /go/src

WORKDIR /go/src

RUN go build -o main .

EXPOSE 7070

RUN chmod 755 main

CMD [ "./main" ]