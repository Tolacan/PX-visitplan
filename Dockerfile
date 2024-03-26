ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

# RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./

COPY events events
COPY repository repository
COPY database database
COPY search search
COPY models models
COPY visitplan-service visitplan-service
COPY query-service query-service

USER root
# Descargar e instalar los paquetes uno por uno
RUN go get -d github.com/elastic/elastic-transport-go/v8@v8.4.0 \
 && go get -d github.com/elastic/go-elasticsearch/v8@v8.12.1 \
 && go get -d github.com/go-logr/logr@v1.3.0 \
 && go get -d github.com/go-logr/stdr@v1.2.2 \
 && go get -d github.com/gorilla/mux@v1.8.1 \
 && go get -d github.com/joho/godotenv@v1.5.1 \
 && go get -d github.com/kelseyhightower/envconfig@v1.4.0 \
 && go get -d github.com/klauspost/compress@v1.17.2 \
 && go get -d github.com/lib/pq@v1.10.9 \
 && go get -d github.com/nats-io/nats.go@v1.33.1 \
 && go get -d github.com/nats-io/nkeys@v0.4.7 \
 && go get -d github.com/nats-io/nuid@v1.0.1 \
 && go get -d github.com/segmentio/ksuid@v1.0.4 \
 && go get -d go.opentelemetry.io/otel@v1.21.0 \
 && go get -d go.opentelemetry.io/otel/metric@v1.21.0 \
 && go get -d go.opentelemetry.io/otel/trace@v1.21.0 \
 && go get -d golang.org/x/crypto@v0.18.0 \
 && go get -d golang.org/x/sys@v0.16.0

# Compilar e instalar la aplicación
RUN go install -a -installsuffix cgo ./...

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=builder /go/bin .