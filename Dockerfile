# Stage 1
FROM golang:1.20-buster AS build
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
RUN apt-get update && apt-get install -y \
  ca-certificates \
  xz-utils
ENV GO111MODULE=on
COPY . .
RUN go mod tidy
RUN go mod download