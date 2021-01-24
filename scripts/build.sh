#! /bin/bash

# install go-bindata
go get -u github.com/go-bindata/go-bindata/...

#embed config.yml
go-bindata -o config.go config

GOARCH=386 go build . 
