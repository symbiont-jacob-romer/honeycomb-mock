FROM golang:1.13.1 as builder

ADD . /src/klyntar
WORKDIR /src/klyntar

# Install binary and strip out unnecessary bytes
ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"
RUN CGO_ENABLED=0 go install --ldflags '-extldflags "-static"' main.go 
RUN strip /go/bin/main

FROM alpine:3.10
COPY --from=builder /go/bin/main .

# --no-cache option allows you to not cache the
RUN apk add --no-cache dumb-init

ENTRYPOINT ["/usr/bin/dumb-init", "./main"]
