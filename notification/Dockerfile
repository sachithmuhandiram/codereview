FROM golang:alpine

RUN apk update
RUN apk add git

RUN go get github.com/sirupsen/logrus

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

RUN go build -o email .

EXPOSE 7072

RUN chmod 755 email

CMD [ "./email" ]