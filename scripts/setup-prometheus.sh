#!/bin/bash
# get prometheus for linux/macos
# wget https://github.com/prometheus/prometheus/releases/download/v2.24.1/prometheus-2.24.1.darwin-amd64.tar.gz
#  tar xvf prometheus-2.24.1.darwin-amd64.tar.gz -C $HOME/apps
"$HOME"/apps/prometheus-2.24.1.darwin-amd64/./prometheus --config.file="$HOME"/apps/prometheus-2.24.1.darwin-amd64/prometheus.yml --log.level=error