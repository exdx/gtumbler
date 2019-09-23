#!/usr/bin/env bash

docker build -t gtumber/client:v0.0.1 Dockerfile.client . && docker build -t gtumbler/mixer:v0.0.1 Dockerfile.mixer