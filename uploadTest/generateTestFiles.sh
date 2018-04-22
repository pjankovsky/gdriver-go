#!/usr/bin/env bash

# make the files in root
for letter in {A..J}; do
    dd if=/dev/zero of=file_${letter} bs=1M count=1 status=none
done

#make flat folder
mkdir folder_flat
for letter in {A..J}; do
    dd if=/dev/zero of=folder_flat/file_${letter} bs=1M count=1 status=none
done

# make recursive folder
for letter in {A..D}; do
    mkdir -p folder_deep/folder_${letter}
    for subletter in {A..D}; do
        dd if=/dev/zero of=folder_deep/folder_${letter}/file_${subletter} bs=1M count=1 status=none
    done
done