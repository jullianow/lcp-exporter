ARG GOLANG_VERSION=1.24

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"

ARG VERSION=latest

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /lcp-exporter

COPY main.go ./
COPY go.* ./
COPY collector ./collector
COPY config ./config
COPY internal ./internal
COPY lcp ./lcp

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.VERSION=${VERSION}" \
    -a -installsuffix cgo \
    -o /go/bin/lcp-exporter \
    ./main.go

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus Exporter for Liferay Cloud Platform (LCP)"
LABEL org.opencontainers.image.source="https://github.com/jullianow/lcp-exporter"

COPY --from=build /go/bin/lcp-exporter /

EXPOSE 9402

ENTRYPOINT ["/lcp-exporte"]
