# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.18-alpine as build

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
ARG TARGETOS TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /order_manager ./cmd/order_manager/main.go 

##
## Deploy
##
FROM scratch

EXPOSE 2223

COPY --from=build /order_manager /

ENTRYPOINT ["/order_manager"]
