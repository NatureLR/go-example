#! /bin/bash

# 修改为自己的仓库地址
image=""

docker build -t $image:latest .
docker push $image:latest

kubectl -n admission-example rollout restart deployment admission-example
kubectl -n admission-example rollout status  deployment admission-example