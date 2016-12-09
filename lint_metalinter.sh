#!/usr/bin/env bash

set +e

gometalinter --deadline=120s --exclude 'should have comment or be unexported' ./...
