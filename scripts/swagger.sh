#!/bin/sh

swag init -g distributionServer.go -d ./pkg/distributionServer,./logstream/delivery/http,./log/delivery/http,./config/delivery/http,./domain,./common/delivery/http
