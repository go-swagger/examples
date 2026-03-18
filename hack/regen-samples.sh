#!/bin/bash
# SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
# SPDX-License-Identifier: Apache-2.0
#
# Thin wrapper: delegates to the Go tool.
# Kept for backward compatibility with CI and existing documentation.

set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "${root}/hack/tools"
exec go run . regen
