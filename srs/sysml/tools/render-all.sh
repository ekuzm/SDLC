#!/bin/sh
# Run from any directory. PNG outputs are committed for offline review.
set -eu
cd "$(dirname "$0")/.."
docker build -t balance-sysml:0.62.0 tools
docker run --rm -u "$(id -u):$(id -g)" -e HOME=/tmp \
  -v "$PWD/..:/srs" -w /srs/sysml balance-sysml:0.62.0 \
  python tools/render.py
