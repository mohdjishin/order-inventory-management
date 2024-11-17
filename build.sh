#!/bin/bash

COMMIT_HASH=$(git rev-parse --short HEAD)
BUILD_TIME=$(date "+%Y-%m-%d %H:%M:%S")

OUTPUT_BIN="oim"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "-X 'github.com/mohdjishin/order-inventory-management/internal/meta.CommitHash=${COMMIT_HASH}' -X 'github.com/mohdjishin/order-inventory-management/internal/meta.BuildTime=${BUILD_TIME}'" -o ${OUTPUT_BIN} ./cmd/api/

if [ $? -eq 0 ]; then
  echo "Build successful: ${OUTPUT_BIN}"
else
  echo "Build failed"
  exit 1
fi
