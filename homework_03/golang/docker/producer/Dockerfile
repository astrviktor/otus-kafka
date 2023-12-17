FROM golang:1.17-buster as build

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/

COPY ./golang/go.mod .
COPY ./golang/go.sum .
RUN go mod download

COPY ./golang/cmd/producer/main.go ./cmd/producer/main.go

RUN go build -o /opt/producer ./cmd/producer/main.go

# ------------------------------------
FROM ubuntu:22.04

LABEL SERVICE="producer"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/producer /srv/producer

WORKDIR /srv
ENTRYPOINT ["/srv/producer"]