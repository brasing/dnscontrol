FROM golang:1.9 AS build-env
WORKDIR /go/src/github.com/StackExchange/dnscontrol
ADD . .
RUN go install .
RUN dnscontrol version

FROM ubuntu:xenial
COPY --from=build-env /go/bin/dnscontrol /
RUN /dnscontrol version
CMD /dnscontrol