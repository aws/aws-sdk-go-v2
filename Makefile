DOC_SRC=$(shell pwd)
SITE_ENV?=production
PREVIEW_HOST?=127.0.0.1
PREVIEW_PORT?=1313

COMMON_DOCKER_BUILD_ARGS=--build-arg SITE_ENV=$(SITE_ENV)

setup: init-submodules init-node-modules

init-node-modules:
	npm install

init-submodules:
	git submodule update --init --recursive --force themes/docsy

generate:
	hugo --environment $(SITE_ENV) -d docs --gc

container-generate:
	docker build \
		$(COMMON_DOCKER_BUILD_ARGS) \
		-t aws-sdk-go-v2-docs \
		-f ./Dockerfile \
		.
	docker run \
		-t aws-sdk-go-v2-docs \
		make setup generate

preview:
	hugo server \
		--bind $(PREVIEW_HOST) \
		--port $(PREVIEW_PORT) \
		--environment $(SITE_ENV) \
		-d docs

container-preview:
	docker build \
		$(COMMON_DOCKER_BUILD_ARGS) \
		-t aws-sdk-go-v2-docs \
		-f ./Dockerfile \
		.
	docker run \
		-p 127.0.0.1:$(PREVIEW_PORT):$(PREVIEW_PORT) \
		--env PREVIEW_HOST=0.0.0.0 \
		--env PREVIEW_PORT=$(PREVIEW_PORT) \
		-i -t aws-sdk-go-v2-docs \
		make preview

.PHONY: setup init-node-modules init-submodules generate container-generate validate preview preview container-preview
