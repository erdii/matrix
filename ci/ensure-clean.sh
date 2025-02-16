#!/bin/bash
set -euxo pipefail

go mod tidy

# Bail if repo is dirty.
test -z "$(git status --porcelain)"
