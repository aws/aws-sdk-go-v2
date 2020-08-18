#!/bin/sh

BASE_PATH=$1
SDK_ROOT=$2
SMITHY_GO_ROOT=$3

for MOD_PATH in $(find ${BASE_PATH} -type f -name "go.mod" | xargs -I {} dirname {})
do
	cd ${MOD_PATH}

	go mod edit --replace github.com/aws/aws-sdk-go-v2=${SDK_ROOT}
	go mod edit --replace github.com/awslabs/smithy-go=${SMITHY_GO_ROOT}
	go mod tidy
done
