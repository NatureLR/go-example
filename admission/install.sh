#! /bin/bash

kubectl apply -f deployment/workload.yaml
kubectl apply -f deployment/mutatingwebhook.yaml
kubectl apply -f deployment/validatingwebhook.yaml