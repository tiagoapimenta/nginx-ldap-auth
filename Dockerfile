FROM golang:alpine

COPY . /go/src/github.com/tiagoapimenta/nginx-ldap-auth

ENV CGO_ENABLED=0

RUN cd /go/src/github.com/tiagoapimenta/nginx-ldap-auth && \
	apk add --no-cache git && \
	go get -u gopkg.in/yaml.v2 && \
	go get -u gopkg.in/ldap.v2 && \
	go build -a -x -ldflags='-s -w -extldflags -static' -v -o /go/bin/nginx-ldap-auth ./main

FROM scratch

MAINTAINER Tiago A. Pimenta <tiagoapimenta@gmail.com>

COPY --from=0 /go/bin/nginx-ldap-auth /usr/local/bin/nginx-ldap-auth

WORKDIR /tmp

VOLUME /etc/nginx-ldap-auth

EXPOSE 5555

USER 65534:65534

CMD [ "/usr/local/bin/nginx-ldap-auth", "--config", "/etc/nginx-ldap-auth/config.yaml" ]
