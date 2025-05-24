#!/usr/bin/env bash
# This file is used to cache the docker in CI

kind load docker-image nginx:1.27 &
kind load docker-image nats:2.10.26 &
kind load docker-image redis:7.4.2 &
# kind load docker-image mysql:8.4.4 &
kind load docker-image bitnami/mysql:8.4.4-debian-12-r10 &
kind load docker-image casbin/casdoor:v1.854.0 &
kind load docker-image python:3.9-slim-buster &
kind load docker-image tdengine/tdengine:3.3.5.2 &
kind load docker-image registry.k8s.io/etcd:3.5.21-0 &
wait
