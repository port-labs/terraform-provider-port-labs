#!/usr/bin/env bash
# Wrapper for gotestsum --raw-command + --rerun-fails.
# Writes -coverprofile only on the first attempt so retries do not overwrite
# the full-package profile with coverage from the failed tests alone.
set -euo pipefail

name=$1
shift

attempt_file="test-results/${name}.attempt"
n=$(cat "$attempt_file")
n=$((n + 1))
echo "$n" > "$attempt_file"

cover_args=()
if [[ "$n" -eq 1 ]]; then
  cover_args=(-coverprofile "coverage/${name}.out")
fi

exec go test -json \
  -timeout 10m \
  -parallel "${TEST_PARALLEL}" \
  -shuffle "${GITHUB_RUN_ID}" \
  "${cover_args[@]}" \
  "$@"
