#!/bin/bash

if [ -z "$1" ]; then
    echo "No coverage threshold argument supplied"
    exit 1
fi
threshold=$1

go test -cover --coverprofile=coverage.out ./...
goc
covered=0
total=0
while IFS='' read -r line || [[ -n "$line" ]]; do
    IFS=' ' read -r -a array <<< "$line"
    total=$(($total+${array[1]}))
    if [ "${array[2]}" != "0" ]; then
        covered=$(($covered+${array[1]}))
    fi
done < "coverage.out"
rm coverage.out
pc=0
if ((${covered} > 0)); then
    pc=$((100*${covered}/${total}))
fi
if (($pc < $threshold)); then
    echo "FAILED: Unit test coverage: $(awk -v covered=$covered -v total=$total 'BEGIN { print (100*covered/total) }')%, threshold: $threshold%, please write more test"
    exit 1
fi
echo "PASS: Unit test coverage: $(awk -v covered=$covered -v total=$total 'BEGIN { print (100*covered/total) }')%, threshold: $threshold%"
