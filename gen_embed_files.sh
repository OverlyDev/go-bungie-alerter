#!/bin/bash

mkdir -p embeds

echo $(date -u) | tr -d '\n' > embeds/time.txt
echo $(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') | tr -d '\n' > embeds/version.txt
echo $(git rev-parse --short HEAD) | tr -d '\n' > embeds/ref.txt