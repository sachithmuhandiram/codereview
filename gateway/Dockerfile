FROM golang:alpine

RUN apk update

COPY . /go/src

WORKDIR /go/src

RUN go build -o main .

EXPOSE 7070

RUN chmod 755 main

CMD [ "./main" ]