#!/bin/bash
set -e

apk update
apk add git curl openssl
apk add py-pip
apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python3-dev make
pip --no-cache-dir install -U pip
