DOC_SRC=$(shell pwd)
SITE_ENV?=production

setup: init-submodules init-node-modules

init-node-modules:
	npm install

init-submodules:
	git submodule update --init --recursive --force themes/docsy

generate:
	hugo --environment $(SITE_ENV) -d docs --gc

.PHONY: setup init-node-modules init-submodules generate
