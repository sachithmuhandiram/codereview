FROM golang:alpine

RUN apk update
RUN apk add git

RUN go get github.com/sirupsen/logrus

# # setting container time zone
# RUN set -e \
#   && apk add tzdata \
#   && cp /usr/share/zoneinfo/Asia/Colombo /etc/localtime \
#   && echo "Asia/Colombo" >  /etc/timezone \
#   && apk del tzdata

COPY . /go/src

WORKDIR /go/src

RUN export LOGINURL=http://localhost:7070/login?token=

RUN go build -o email .

EXPOSE 7072

RUN chmod 755 email

CMD [ "./email" ]