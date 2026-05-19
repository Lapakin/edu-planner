#!/bin/bash

## grep is used to exclude packages from testing, if you want to add
## a new one excluded package you need to do that:
## go test $(go list ./... | grep -v *your package*) | grep -v '\[no test files\]'

INTEGRATION=false go test $(go list ./...) | grep -v '\[no test files\]'
exit ${PIPESTATUS[0]}
