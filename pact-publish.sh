#!/bin/bash

set -x

VERSION=0.0.1 #like 1.0.0
BROKER_BASE_URL=http://localhost
TAG=dev
BRANCH=basket-service-pact
BASKET_SERVICE_PRODUCT_SERVICE_PACT=./app/product/pacts/basketservice-productservice.json

pact-broker publish \
${BASKET_SERVICE_PRODUCT_SERVICE_PACT} \
--consumer-app-version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL} \
--tag=${TAG} \
--branch=${BRANCH}
