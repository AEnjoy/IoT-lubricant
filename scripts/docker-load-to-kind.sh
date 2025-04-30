#!/usr/bin/env bash
# This file is used to cache the docker in CI
kind load docker-image amd64/mongo:8.0-noble &
kind load docker-image docker.elastic.co/elasticsearch/elasticsearch:8.13 &
kind load docker-image graylog/graylog:6.2 &
kind load docker-image nginx:1.27 &
kind load docker-image nats:2.10.26 &
kind load docker-image redis:7.4.2 &
kind load docker-image bitnami/mysql:8.4.4-debian-12-r4 &
kind load docker-image casbin/casdoor:v1.854.0 &
kind load docker-image python:3.9-slim-buster &
wait
