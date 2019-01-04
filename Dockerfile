FROM golang:1.10-alpine3.7 as build

RUN apk add -U git make

WORKDIR /go/src/github.com/keycloak/keycloak-gatekeeper
COPY . .
RUN make static

FROM alpine:3.7

RUN apk add --no-cache ca-certificates

ADD templates/ /opt/templates
COPY --from=build /go/src/github.com/keycloak/keycloak-gatekeeper/bin/keycloak-gatekeeper /opt/keycloak-gatekeeper

WORKDIR "/opt"

ENTRYPOINT [ "/opt/keycloak-gatekeeper" ]
