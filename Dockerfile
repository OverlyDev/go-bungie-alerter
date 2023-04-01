# Grabber + packer stage
FROM alpine as grabpack

RUN apk add --no-cache upx wget

WORKDIR /grabpack

RUN arch=$(arch | sed s/aarch64/linux-arm64/ | sed s/x86_64/linux-amd64/) && \
    wget https://github.com/OverlyDev/go-bungie-alerter/releases/latest/download/BungieAlerter-${arch} && \
    chmod +x BungieAlerter-${arch} && \
    upx --best --lzma -o /grabpack/BungieAlerter /grabpack/BungieAlerter-${arch} && \
    chmod +x BungieAlerter
    
# Final stage
FROM alpine
COPY --from=grabpack /grabpack/BungieAlerter /app/BungieAlerter
CMD ["/app/BungieAlerter"]