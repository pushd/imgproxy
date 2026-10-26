#!/bin/bash

set -e

go build -v -ldflags "-s -w" -o /usr/local/bin/imgproxy
