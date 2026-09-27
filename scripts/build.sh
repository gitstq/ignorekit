#!/usr/bin/env bash
# Cross-platform build for ignorekit.
# Produces archives + SHA256 checksums under dist/.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="${VERSION:-1.0.0}"
COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

DIST="$ROOT/dist"
rm -rf "$DIST"
mkdir -p "$DIST"

# GOOS/GOARCH targets: Windows, macOS, Linux (amd64 + arm64; plus 386).
TARGETS=(
  "windows amd64"
  "windows arm64"
  "windows 386"
  "darwin  amd64"
  "darwin  arm64"
  "linux   amd64"
  "linux   arm64"
  "linux   386"
)

LDFLAGS="-s -w \
  -X github.com/gitstq/ignorekit/cmd.Version=${VERSION} \
  -X github.com/gitstq/ignorekit/cmd.Commit=${COMMIT} \
  -X github.com/gitstq/ignorekit/cmd.BuildDate=${DATE}"

for t in "${TARGETS[@]}"; do
  GOOS="${t%% *}"
  GOARCH="${t##* }"
  NAME="ignorekit_${VERSION}_${GOOS}_${GOARCH}"
  WORK="$DIST/$NAME"
  mkdir -p "$WORK"

  BIN="ignorekit"
  [ "$GOOS" = "windows" ] && BIN="ignorekit.exe"

  echo ">> building $GOOS/$GOARCH"
  CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -trimpath -ldflags "$LDFLAGS" -o "$WORK/$BIN" .

  cp LICENSE "$WORK/"
  cp README.md "$WORK/" 2>/dev/null || true
  mkdir -p "$WORK/docs"
  cp docs/README.zh-CN.md docs/README.zh-TW.md "$WORK/docs/" 2>/dev/null || true
  cp CONTRIBUTING.md "$WORK/" 2>/dev/null || true

  if [ "$GOOS" = "windows" ]; then
    (cd "$DIST" && zip -q -r "$NAME.zip" "$NAME")
  else
    (cd "$DIST" && tar -czf "$NAME.tar.gz" "$NAME")
  fi
  rm -rf "$WORK"
done

echo ">> generating SHA256SUMS"
(cd "$DIST" && find . -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) -printf '%f\n' \
  | sort | xargs sha256sum) > "$DIST/SHA256SUMS"

echo ">> done"
ls -1 "$DIST"
