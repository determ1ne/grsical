#!/bin/sh

set -e

GRSICAL_BIN="./grsical-linux-"

is_darwin() {
  case "$(uname -s)" in
  *darwin* ) true ;;
  *Darwin* ) true ;;
  * ) false;;
  esac
}

if is_darwin; then
  zsh ./launch.command
  exit 0
fi

case "$(uname -m)" in
x86_64 ) GRSICAL_BIN="${GRSICAL_BIN}amd64";;
aarch64 ) GRSICAL_BIN="${GRSICAL_BIN}arm64";;
* ) echo "unsupported architecture $(uname -m)}"; exit 1;;
esac

$GRSICAL_BIN -i upfile.json -c config.json -t tweaks.json
