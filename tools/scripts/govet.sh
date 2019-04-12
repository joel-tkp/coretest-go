#!/bin/bash

# If any checks are found to be useless, they can be disabled here.
# See the output of "go tool vet" for a list of flags.
vetflags="-all=true"

errors=

if ! go vet ./... 2>&1; then
	errors=YES
fi

[ -z  "$errors" ] && exit 0

echo
echo "Please fix the go vet warnings above. To disable certain checks, change vetflags."
exit 1