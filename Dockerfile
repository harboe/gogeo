FROM golang:latest
MAINTAINER dennis.harboe@gmail.com

ADD . /go/src/github.com/harboe/gogeo
WORKDIR /go/src/github.com/harboe/gogeo

# building the application
RUN go install

EXPOSE :8080
ENTRYPOINT ["gogeo"]

# execute the app
CMD ["http", "-p", ":8080"]
