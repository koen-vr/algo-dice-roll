#!/bin/bash

date '+keyreg-teal-test start %Y%m%d_%H%M%S'

set -e
set -x
set -o pipefail
export SHELLOPTS

PYTHON=python3
gcmd="goal -d ../../net1/primary"
ACC1=$(${gcmd} account list | awk '{ print $3 }' | head -n 1)

PYTEAL_APPROVAL_PROG="../../contracts/roller.py"
PYTEAL_CLEAR_PROG="../../contracts/clear.py"
TEAL_APPROVAL_PROG="../../contracts/roller.teal"
TEAL_CLEAR_PROG="../../contracts/clear.teal"

# compile PyTeal into TEAL
"$PYTHON" "$PYTEAL_APPROVAL_PROG" > "$TEAL_APPROVAL_PROG"
"$PYTHON" "$PYTEAL_CLEAR_PROG" > "$TEAL_CLEAR_PROG"

# create app
APP_ID=$(
  ${gcmd} app create --creator "$ACC1" \
    --approval-prog "$TEAL_APPROVAL_PROG" \
    --clear-prog "$TEAL_CLEAR_PROG" \
    --global-byteslices 1 \
    --global-ints 1 \
    --local-byteslices 0 \
    --local-ints 0 \
    --on-completion OptIn |
    grep Created |
    awk '{ print $6 }'
)
echo "App ID = $APP_ID"

# read global state
STATE=$(${gcmd} app read --app-id "$APP_ID" --guess-format --local --from "$ACC1")
STATE=$(${gcmd} app read --app-id "$APP_ID" --guess-format --global)
