#!/bin/sh

SDK_ROOT=$1
CODGEN_ROOT=$2
SDK_PREFIX=$3

for MOD_PATH in $(find ${CODGEN_ROOT} -type f -name "go.mod" | xargs -I {} dirname {})
do
	cd ${MOD_PATH}
	MOD=$(go list -m)

	DST=${SDK_ROOT}/${MOD#"$SDK_PREFIX"}

	echo "Copying ${MOD} to ${DST}"

	rm -rf ${DST}
	mkdir -p ${DST}
	cp -r . ${DST}
done
