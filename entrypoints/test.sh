#!/usr/bin/env bash
set -e

ginkgo -r --randomizeAllSpecs --randomizeSuites --trace --progress --succinct .

# then watch if asked to
if [[ $1 = "watch" ]] || [[ $1 = "-w" ]]; then
  ginkgo watch -v -r --randomizeAllSpecs --failOnPending --trace --progress --succinct .
fi

exit $?
