#!/bin/bash

set -e

OWL_SECRET_KEY=test-secret-key \
go test -v -coverprofile=coverage.out ./...
