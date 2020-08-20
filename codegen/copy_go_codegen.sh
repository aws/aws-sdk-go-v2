#!/bin/sh

SDK_ROOT=$1
CODGEN_ROOT=$2
SDK_PREFIX=$3

for MOD_PATH in $(find ${CODGEN_ROOT} -type f -name "go.mod" | xargs -I {} dirname {})
do
	cd ${MOD_PATH}
	MOD=$(grep -e "^module" go.mod | sed -e 's/^module *//')

	DST=${SDK_ROOT}/${MOD#"$SDK_PREFIX"}

	echo "Copying ${MOD} to ${DST}"

	rm ${DST}/*.go
	rm -rf ${DST}/types
	rm -rf ${DST}/internal/endpoints
	mkdir -p ${DST} 2>/dev/null
	cp -r . ${DST}
done
