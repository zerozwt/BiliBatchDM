#!/usr/bin/env bash

cd $(dirname "$0")

rm -rf output
mkdir output

cd server
go build -o ../output/BiliBatchDM
cd ..

cd www
npm run build
mv dist ../output/
cd ..