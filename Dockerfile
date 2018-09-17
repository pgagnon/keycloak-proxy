FROM golang:1.10-alpine3.7 as build

RUN apk add -U git make

WORKDIR /go/src/github.com/gambol99/keycloak-proxy
COPY . .
RUN make static

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

ADD templates/ /opt/templates
COPY --from=build /go/src/github.com/gambol99/keycloak-proxy/bin/keycloak-proxy /opt/keycloak-proxy

WORKDIR "/opt"

ENTRYPOINT [ "/opt/keycloak-proxy" ]
