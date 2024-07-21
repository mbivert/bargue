#!/bin/sh

# Check that dependencies (executables) are installed
# Fails (exit 1) if not.

set -e

which go      > /dev/null
# See https://github.com/mbivert/dtmpl/
which dtmpl   > /dev/null
