# Builder stage
FROM golang:alpine AS builder

RUN apk add --no-cache bash git

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.sh ./
RUN go generate && go build -ldflags "-s -w" -o BA-base

# Packer stage
FROM alpine as packer
COPY --from=builder /build/BA-base /packer/BA-base
RUN apk add upx && upx --best --lzma -o /packer/BA-packed /packer/BA-base

# Final stage
FROM alpine
# COPY --from=builder /build/hello /app/hello
COPY --from=packer /packer/BA-packed /app/BungieAlerter
CMD ["/app/BungieAlerter"]