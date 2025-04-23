#!/usr/bin/env bash
# This file is used to cache the docker in CI

echo "Cache: Saving Docker Images..."
mkdir -p /tmp/docker-images
docker save \
    nginx:1.27 alpine:3.21.3 kindest/node:v1.32.2 \
    nats:2.10.26 redis:7.4.2 casbin/casdoor:v1.854.0 \
    python:3.9-slim-buster golang:1.24-alpine \
    bitnami/mysql:8.4.4-debian-12-r10 \
    tdengine/tdengine:3.3.6.3 \
    registry.k8s.io/etcd:3.5.21-0 \
    -o /tmp/docker-images/base-images.tar &
