FROM golang:1.25.5-alpine3.23 AS builder
ARG TARGET

WORKDIR /src

RUN mkdir ./bin

COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    --mount=type=cache,id=gocache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    --mount=type=cache,id=gocache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GODEBUG=x509usefallbackroots=1 \
    go build -trimpath -tags timetzdata -o ./bin/ ./cmd/${TARGET}

FROM scratch
ARG TARGET

COPY --from=builder /src/bin/${TARGET} /app/bin/app

USER 1000:1000

WORKDIR /app

ENTRYPOINT ["/app/bin/app"]
