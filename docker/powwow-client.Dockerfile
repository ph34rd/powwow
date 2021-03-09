FROM golang:latest AS builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /build
COPY . .
RUN make powwow-client

FROM alpine:latest AS single
RUN apk add --no-cache ca-certificates
RUN test -e /etc/nsswitch.conf || echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder ./build/build/bin/powwow-client .
ENV POWWOW_HOST localhost
ENV POWWOW_PORT 9999
CMD ./powwow-client -addr "${POWWOW_HOST}:${POWWOW_PORT}"

FROM alpine:latest AS parallel
RUN apk add --no-cache ca-certificates parallel
RUN test -e /etc/nsswitch.conf || echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=builder ./build/build/bin/powwow-client .
ENV POWWOW_HOST localhost
ENV POWWOW_PORT 9999
ENV POWWOW_PARALLEL 10
CMD seq 1 ${POWWOW_PARALLEL} | parallel --will-cite -j 0 -I{} ./powwow-client -addr "${POWWOW_HOST}:${POWWOW_PORT}"
