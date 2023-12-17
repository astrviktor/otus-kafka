FROM golang:1.17-buster as build

ENV BIN_FILE "/opt/consumer/consumer"
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/

COPY ./golang/go.mod .
COPY ./golang/go.sum .
RUN go mod download

COPY ./golang/cmd/consumer/main.go ./cmd/consumer/main.go

ARG LDFLAGS
#RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o /opt/producer ./cmd/producer/main.go
RUN go build -o ${BIN_FILE} ./cmd/consumer/main.go

# ------------------------------------
FROM ubuntu:22.04

LABEL SERVICE="consumer"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/consumer/consumer /srv/consumer

WORKDIR /srv
ENTRYPOINT ["/srv/consumer"]
