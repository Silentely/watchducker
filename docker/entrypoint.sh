#!/bin/sh
set -e

TZ_VALUE="${TZ:-UTC}"
ZONEINFO_PATH="/usr/share/zoneinfo/${TZ_VALUE}"

if [ -f "$ZONEINFO_PATH" ]; then
  ln -snf "$ZONEINFO_PATH" /etc/localtime
  echo "$TZ_VALUE" >/etc/timezone
else
  echo "Warning: timezone '${TZ_VALUE}' not found, falling back to UTC" >&2
  ln -snf /usr/share/zoneinfo/UTC /etc/localtime
  echo "UTC" >/etc/timezone
fi

exec "$@"
