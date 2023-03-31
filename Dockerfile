# Builder stage
FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.sh ./
RUN go build -ldflags "-s -w" -o hello

# Packer stage
FROM alpine as packer
COPY --from=builder /build/hello /packer/hello-stripped
RUN apk add upx && upx --best --lzma -o /packer/hello /packer/hello-stripped

# Final stage
FROM alpine
# COPY --from=builder /build/hello /app/hello
COPY --from=packer /packer/hello /app/BungieAlerter
CMD ["/app/BungieAlerter"]