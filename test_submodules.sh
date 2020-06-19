#!/bin/sh

BASE_PATH=$1
TEST_CMD=$2

for MOD_PATH in $(find ${BASE_PATH} -type f -name "go.mod" | xargs -I {} dirname {})
do
	cd ${MOD_PATH}
	MOD=$(go list -m)
	echo "Testing ${MOD}"

	${TEST_CMD}
done

