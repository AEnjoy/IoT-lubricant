#!/usr/bin/env bash
# This file is used to cache the docker in CI

kind load docker-image nginx:1.27 &
kind load docker-image nats:2.10.26 &
kind load docker-image redis:7.4.2 &
kind load docker-image bitnami/mysql:8.4.4-debian-12-r4 &
kind load docker-image casbin/casdoor:v1.854.0 &
kind load docker-image python:3.9-slim-buster &
wait
