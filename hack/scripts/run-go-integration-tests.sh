#!/bin/bash

## grep is used to exclude packages from testing, if you want to add
## a new one exluded package you need to do that:
## go test $(go list ./... | grep -v *your package*) | grep -v '\[no test files\]'

if INTEGRATION=true go test $(go list ./...) | grep -v '\[no test files\]' | grep FAIL; then
  exit 1  # At least one test failed
else
  exit 0  # All tests passed
fi