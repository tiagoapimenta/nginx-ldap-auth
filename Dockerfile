FROM golang:1.13-alpine AS gobuild

COPY . /build/nginx-ldap-auth

ENV CGO_ENABLED=0

RUN cd /build/nginx-ldap-auth && \
	apk add --no-cache git && \
	go build -a -x -ldflags='-s -w -extldflags -static' -v -o /go/bin/nginx-ldap-auth ./main

FROM scratch

MAINTAINER Tiago A. Pimenta <tiagoapimenta@gmail.com>

COPY --from=gobuild /go/bin/nginx-ldap-auth /usr/local/bin/nginx-ldap-auth

WORKDIR /tmp

VOLUME /etc/nginx-ldap-auth

EXPOSE 5555

USER 65534:65534

CMD [ \
	"/usr/local/bin/nginx-ldap-auth", \
	"--config", \
	"/etc/nginx-ldap-auth/config.yaml" \
]
