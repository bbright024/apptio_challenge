From golang:1.10.1-alpine
MAINTAINER Brian Bright bbright123@yahoo.com

ENV SOURCES /go/src/apptio/

COPY . ${SOURCES}

RUN cd ${SOURCES}/logserver && GOBIN="/bin/" go install && cd / && touch logserver.log

#touch /logserver.log && cp /go/src/apptio/logserver/docker_deploy_conf.json /

EXPOSE 8888

ENTRYPOINT logserver ${SOURCES}/logserver/docker_deploy_conf.json

