ARG GOLANG_VERSION=1.24

ARG TARGETOS="linux"
ARG TARGETARCH="amd64"

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /lcp-exporter

COPY main.go ./
COPY go.* ./
COPY collector ./collector
COPY config ./config
COPY internal ./internal
COPY lcp ./lcp

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.VERSION=${VERSION}" \
    -a -installsuffix cgo \
    -o /go/bin/lcp-exporte \
    ./main.go

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus Exporter for Liferay Cloud Platform (LCP)"
LABEL org.opencontainers.image.source="https://github.com/jullianow/lcp-exporte"

COPY --from=build /go/bin/lcp-exporte /

EXPOSE 9402

ENTRYPOINT ["/lcp-exporte"]
