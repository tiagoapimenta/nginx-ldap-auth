FROM golang:alpine

COPY . /go/src/github.com/tiagoapimenta/nginx-ldap-auth

RUN cd /go/src/github.com/tiagoapimenta/nginx-ldap-auth && \
	apk add --no-cache git && \
	go get -u gopkg.in/yaml.v2 && \
	go get -u gopkg.in/ldap.v2 && \
	go build -ldflags='-s -w' -v -o /go/bin/nginx-ldap-auth ./main

FROM alpine

MAINTAINER Tiago A. Pimenta <tiagoapimenta@gmail.com>

COPY --from=0 /go/bin/nginx-ldap-auth /usr/local/bin/nginx-ldap-auth

WORKDIR /tmp

VOLUME /etc/nginx-ldap-auth

EXPOSE 5555

USER nobody

CMD [ "nginx-ldap-auth", "--config", "/etc/nginx-ldap-auth/config.yaml" ]
