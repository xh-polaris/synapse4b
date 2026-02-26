#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=synapse4b
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}