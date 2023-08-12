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
RUN make build

# Stage 2: Minimize the binary using https://github.com/upx/upx
FROM alpine:latest AS compress
WORKDIR /
RUN apk update && \
  apk add xz
ADD https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.96-amd64_linux.tar.xz | \
  tar -xOf - upx-3.96-amd64_linux/upx > /bin/upx && \
  chmod a+x /bin/upx
COPY --from=build /go/src/app/build/teacher-api /teacher-api
RUN upx ./teacher-api

# Stage 3
FROM scratch
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/.env /.env
COPY --from=compress /teacher-api /teacher-api
EXPOSE 3000
CMD ["./teacher-api"]