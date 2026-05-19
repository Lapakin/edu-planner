#!/bin/bash

## grep is used to exclude packages from testing, if you want to add
## a new one exluded package you need to do that:
## go test $(go list ./... | grep -v *your package*) | grep -v '\[no test files\]'

INTEGRATION=true go test -timeout 300s $(go list ./...) | grep -v '\[no test files\]'
exit ${PIPESTATUS[0]}
