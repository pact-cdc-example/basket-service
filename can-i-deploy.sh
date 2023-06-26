#!/bin/bash

set -x

VERSION=0.0.1
PACTICIPANT=BasketService
BROKER_BASE_URL=http://localhost
ENV=dev
#
#pact-broker can-i-deploy \
#--pacticipant=${PACTICIPANT} \
#--broker-base-url=${BROKER_BASE_URL} \
#--latest

pact-broker can-i-deploy --pacticipant ${PACTICIPANT} --version ${VERSION} --to dev --retry-while-unknown=12 --retry-interval=10 --limit 1000