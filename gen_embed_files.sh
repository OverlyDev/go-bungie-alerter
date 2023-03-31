#!/bin/bash

mkdir -p embeds

echo $(date -u) > embeds/time.txt
echo $(git describe --tags --abbrev=0 2>/dev/null || echo 'dev') > embeds/version.txt
echo $(git rev-parse --short HEAD) > embeds/ref.txt