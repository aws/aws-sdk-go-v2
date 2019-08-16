#!/bin/sh

SRC_REPO=https://github.com/aws/aws-sdk-go
SRC_ROOT=/tmp/aws-sdk-go-models

SRC_APIS_ROOT=${SRC_ROOT}/models/apis
DST_APIS_ROOT=models/apis

SRC_ENDPOINTS=${SRC_ROOT}/models/endpoints/endpoints.json
DST_ENDPOINTS=models/endpoints/endpoints.json

rm -rf ${SRC_ROOT}
git clone --single-branch --branch master ${SRC_REPO} ${SRC_ROOT}

rm -rf ${DST_APIS_ROOT}/*
rm -rf ${DST_ENDPOINTS}

cp -r ${SRC_APIS_ROOT}/* ${DST_APIS_ROOT}/
cp ${SRC_ENDPOINTS} ${DST_ENDPOINTS}

make generate
