FROM golang:1.27-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=$GOPATH/pkg/mod \
	go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=$GOPATH/pkg/mod \
	CGO_ENABLED=0 \
	go build -ldflags="-w -s" -trimpath \
	-o ./bot ./cmd/app


FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates \
	&& adduser -D -H appuser
COPY --chown=appuser:appuser --from=builder /app/bot .

USER appuser

ENTRYPOINT [ "./bot" ]
