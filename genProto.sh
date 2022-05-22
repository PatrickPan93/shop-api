#!/bin/zsh
protoc -I ./user-web/proto user.proto --go_out=plugins=grpc:.
