#!/bin/bash

docker run --rm -e GOBIN=/go/bin/ -v "$PWD"/bin:/go/bin/ -v "$PWD"/src/:/go/src/ -w /go/src/ golang go install ./roomserver
