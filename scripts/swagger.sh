#!/bin/sh

swag init -g distribution_server.go -d ./pkg/distributionServer,./logstream/delivery/http,./log/delivery/http,./config/delivery/http,./domain,./common/delivery/http
