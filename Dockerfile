FROM golang:1.26.4-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /balance ./cmd/balance
FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 balance && mkdir /data && chown balance /data
COPY --from=build /balance /usr/local/bin/balance
USER balance
ENV BALANCE_ADDR=0.0.0.0:8080 BALANCE_DATA=/data/balance.json
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["balance"]
