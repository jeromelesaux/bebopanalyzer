#!/bin/bash

GOOS=linux GOARCH=amd64 make package

GOOS=linux GOARCH=ppc64 make package

GOOS=darwin GOARCH=amd64 make package

GOOS=darwin GOARCH=386 make package

GOOS=windows GOARCH=amd64 make package

GOOS=windows GOARCH=386 make package

GOOS=linux GOARCH=arm make package