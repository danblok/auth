# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine as base
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

FROM base as builder
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/server cmd/server/main.go

FROM scratch as server
COPY --from=builder /bin/server /bin/
EXPOSE 3000
CMD [ "/bin/server" ]
