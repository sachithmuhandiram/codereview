FROM golang:alpine

RUN apk update

COPY . /go/src

WORKDIR /go/src

RUN go build -o email .

EXPOSE 7071

RUN chmod 755 email

CMD [ "./email" ]