#!/usr/bin/env bash
set -Eeuo pipefail

# From https://github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/dist/install.sh
#
# Installation script for textile. It tries to move $bin in one of the
# directories stored in $binpaths.

INSTALL_DIR="$(dirname "$0")"

bin="$INSTALL_DIR/textile"
binpaths="/usr/local/bin /usr/bin"

# This variable contains a nonzero length string in case the script fails
# because of missing write permissions.
is_write_perm_missing=""

for binpath in $binpaths; do
	if mv "$bin" "$binpath/$bin" 2>/dev/null; then
		echo "Moved $bin to $binpath"
		exit 0
	else
		if test -d "$binpath" && ! test -w "$binpath"; then
			is_write_perm_missing=1
		fi
	fi
done

echo "We cannot install $bin in one of the directories $binpaths"

if test -n "$is_write_perm_missing"; then
	echo "It seems that we do not have the necessary write permissions."
	echo "Perhaps try running this script as a privileged user:"
	echo
	echo "    sudo $0"
	echo
fi

exit 1
