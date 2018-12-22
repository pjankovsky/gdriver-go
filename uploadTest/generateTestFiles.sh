#!/usr/bin/env bash

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    DDCMD="dd"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    command -v gdd >/dev/null 2>&1 || { echo >&2 "I require 'gdd' but it's not installed. Install with 'brew install coreutils'."; exit 1; }
    DDCMD="gdd"
else
    DDCMD="dd"
fi

# make the files in root
for letter in {A..J}; do
    ${DDCMD} if=/dev/zero of=file_${letter} bs=1M count=1 status=none
done

#make flat folder
mkdir -p folder_flat
for letter in {A..J}; do
    ${DDCMD} if=/dev/zero of=folder_flat/file_${letter} bs=1M count=1 status=none
done

# make recursive folder
for letter in {A..D}; do
    mkdir -p folder_deep/folder_${letter}
    for subletter in {A..D}; do
        ${DDCMD} if=/dev/zero of=folder_deep/folder_${letter}/file_${subletter} bs=1M count=1 status=none
    done
done