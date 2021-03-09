FROM golang:latest AS builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /build
COPY . .
RUN make powwow-server

FROM alpine:latest AS production
RUN apk add --no-cache ca-certificates
RUN test -e /etc/nsswitch.conf || echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder ./build/build/bin/powwow-server .
ENV POWWOW_PORT 9999
CMD ./powwow-server -bind "${POWWOW_HOST}:${POWWOW_PORT}"

FROM alpine:latest AS development
RUN apk add --no-cache ca-certificates
RUN test -e /etc/nsswitch.conf || echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder ./build/build/bin/powwow-server .
ENV POWWOW_PORT 9999
CMD ./powwow-server -bind "${POWWOW_HOST}:${POWWOW_PORT}" -dev
