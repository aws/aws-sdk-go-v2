DOC_SRC=$(shell pwd)
SITE_ENV?=production

setup: init-submodules init-node-modules apply-patches

init-node-modules:
	npm install

init-submodules:
	git submodule update --init --recursive --force themes/docsy

apply-patches:
	cd themes/docsy && \
		git apply $(DOC_SRC)/patches/docsy/0001-Update-Bootstrap-Version-to-4.5.3.patch

generate:
	hugo --environment $(SITE_ENV) -d docs --gc

.PHONY: setup init-node-modules init-submodules apply-patches generate
