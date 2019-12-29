ARG ALPINE_VERSION=3.11
ARG GOLANG_VERSION=1.13.5

# FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS modd
# RUN wget https://github.com/cortesi/modd/releases/download/v0.8/modd-0.8-linux64.tgz && \
#     tar -xzf modd-0.8-linux64.tgz && \
#     cp modd-0.8-linux64/modd /modd

FROM alpine:${ALPINE_VERSION} AS dockerize
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS base
RUN mkdir /tmp/unprivileged && \
    echo 'service:x:101:101:service:/:' > /tmp/unprivileged/passwd && \
    echo 'service:x:101:' > /tmp/unprivileged/group
RUN apk update && \
    apk add --no-cache ca-certificates git tzdata && \
    update-ca-certificates
WORKDIR /go/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

FROM base AS builder
RUN CGO_ENABLED=0 go build -ldflags '-d' -tags netgo -installsuffix netgo -o /backend-api cmd/backend/main.go

FROM alpine:${ALPINE_VERSION}
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin/dockerize
COPY --from=base /tmp/unprivileged/group /tmp/unprivileged/passwd /etc/
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /go/src/app/backend/migrations /app/backend/migrations
COPY --from=base /go/src/app/static /app/static
COPY --from=builder /backend-api /app/backend-api
USER 101:101
WORKDIR /app
ENTRYPOINT ["dockerize", "-wait", "tcp://pg:5432", "./backend-api"]
