# syntax=docker/dockerfile:1

#build binary
FROM golang:1.20-alpine AS build-stage
WORKDIR /backend
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver # CGO_ENABLED: makes the binary statically-linked, GOOS: specifies the OS to build for

#minimize image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
# bare minimum distro without any extra components
WORKDIR /
COPY --from=build-stage /webserver /webserver
EXPOSE 80
ENTRYPOINT ["/webserver"]
