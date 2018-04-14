From golang:1.10.1-alpine
MAINTAINER bbright123


COPY . /go/src/apptio

RUN cd src/apptio/ && GOBIN="/bin/" go install /go/src/apptio/logserver/logserver.go 

EXPOSE 8888

ENTRYPOINT logserver /go/src/apptio/logserver/conf.json

