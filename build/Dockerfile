FROM golang:1.21.4-alpine AS builder

RUN /sbin/apk update && \
	/sbin/apk --no-cache add ca-certificates git tzdata && \
	/usr/sbin/update-ca-certificates

RUN adduser -D -g '' api
WORKDIR /home/api

COPY go.mod      ./
COPY go.sum      ./
COPY cmd         ./cmd/
COPY internal    ./internal/
# COPY tls        ./tls/
# COPY ui         ./ui/
COPY *.go        ./
COPY config.yaml ./
# RUN chmod 644  ./tls/key.pem

ARG VERSION

RUN CGO_ENABLED=0 go build -a -tags netgo,osusergo \
    -ldflags "-extldflags '-static' -s -w" \
    -ldflags "-X main.version=${VERSION}" -o api ./cmd/api

FROM busybox:musl

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /home/api/api /home/api
# COPY --from=builder /home/web/tls /home/tls/

USER api
WORKDIR /home
EXPOSE 8080

ENV MYSQL_PASSWORD=$MYSQL_PASSWORD

ENTRYPOINT ["./api"]
CMD ["-addr=:8080", "-debug=false"]