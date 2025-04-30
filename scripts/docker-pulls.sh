#!/usr/bin/env bash
# This file is used to cache the docker in CI

# Docker Pull List
# for base images
docker pull golang:1.24-alpine &
docker pull alpine:3.21.3 &
docker pull kindest/node:v1.32.2 &
docker pull python:3.9-slim-buster &
# for Kubernetes yaml
docker pull amd64/mongo:8.0-noble &
docker pull docker.elastic.co/elasticsearch/elasticsearch:8.13 &
docker pull graylog/graylog:6.2 &
docker pull nginx:1.27 &
docker pull nats:2.10.26 &
docker pull redis:7.4.2 &
docker pull bitnami/mysql:8.4.4-debian-12-r4 &
docker pull casbin/casdoor:v1.854.0 &
wait
