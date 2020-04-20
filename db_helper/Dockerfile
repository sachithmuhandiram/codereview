FROM golang:alpine

RUN apk update
RUN apk add git

RUN set -e \
    && go get github.com/sirupsen/logrus \
    && go get github.com/carlescere/scheduler \
    && go get github.com/go-sql-driver/mysql

# # setting container time zone
RUN \
    apk --update add curl bash nano tzdata && \
    cp /usr/share/zoneinfo/Asia/Colombo /etc/localtime && \
    echo "Asia/Colombo" > /etc/timezone && \
    apk del tzdata && \
    rm -r /var/cache/apk/* && \
    mkdir -p /usr/share/zoneinfo/Asia && \
    ln -s /etc/localtime /usr/share/zoneinfo/Asia/Colombo

COPY . /go/src

WORKDIR /go/src

RUN go build -o main .

EXPOSE 7073

RUN chmod 755 main

CMD [ "./main" ]