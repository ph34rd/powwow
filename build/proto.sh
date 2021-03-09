#!/usr/bin/env bash

protoc --proto_path="$(go list -f '{{ .Dir }}' -m github.com/gogo/protobuf)" -I. --gogofaster_out=. "*.proto"
