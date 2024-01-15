FROM golang:1.19-buster as build

ENV BIN_FILE "/opt/processor"
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/

COPY go.mod .
COPY go.sum .
RUN ls
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN ls

RUN go build -o ${BIN_FILE} ./cmd/processor/main.go

# ------------------------------------
FROM ubuntu:22.04

LABEL SERVICE="processor"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/processor /srv/processor

WORKDIR /srv
ENTRYPOINT ["/srv/processor"]
