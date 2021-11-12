#!/bin/bash
set -e

apt update
apt install -y git curl openssl
apt install -y pip
apt install -y gcc libffi-dev libssl-dev python3-dev make unzip wget
#apt install -y --virtual=build gcc libffi-dev musl-dev openssl-dev python3-dev make
pip --no-cache-dir install -U pip
