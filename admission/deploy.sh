#! /bin/bash

set -e

# 修改为自己的仓库地址
image="naturelingran"

docker buildx build --platform linux/amd64,linux/arm64 -t $image/admission-example:latest -o type=registry .

kubectl -n admission-example rollout restart deployment admission-example
kubectl -n admission-example rollout status  deployment admission-example