# Lint rules to ignore
LINTIGNORESINGLEFIGHT='internal/sync/singleflight/singleflight.go:.+error should be the last type'

UNIT_TEST_TAGS=
BUILD_TAGS=-tags "example,codegen,integration,ec2env"

EACHMODULE_FAILFAST ?= true
EACHMODULE_FAILFAST_FLAG=-fail-fast=${EACHMODULE_FAILFAST}

EACHMODULE_CONCURRENCY ?= 1
EACHMODULE_CONCURRENCY_FLAG=-c ${EACHMODULE_CONCURRENCY}

# SDK's Core and client packages that are compatable with Go 1.9+.
SDK_CORE_PKGS=./aws/... ./internal/...
SDK_CLIENT_PKGS=./service/...
SDK_COMPA_PKGS=${SDK_CORE_PKGS} ${SDK_CLIENT_PKGS}

# SDK additional packages that are used for development of the SDK.
SDK_EXAMPLES_PKGS=
SDK_ALL_PKGS=${SDK_COMPA_PKGS} ${SDK_EXAMPLES_PKGS}

RUN_NONE=-run '^$$'
RUN_INTEG=-run '^TestInteg_'

LICENSE_FILE=$(shell pwd)/LICENSE.txt

all: generate unit

###################
# Code Generation #
###################
generate: smithy-generate gen-config-asserts gen-repo-mod-replace tidy-modules-. add-module-license-files gen-aws-ptrs

smithy-generate:
	cd codegen && ./gradlew clean build -Plog-tests && ./gradlew clean

gen-config-asserts:
	@echo "Generating SDK config package implementor assertions"
	cd config && go generate

gen-repo-mod-replace:
	@echo "Generating go.mod replace for repo modules"
	cd internal/repotools/cmd/makerelative && go run ./

gen-mod-replace-local:
	./mod_replace_local_submodules.sh `pwd` `pwd` `pwd`/../smithy-go

gen-aws-ptrs:
	cd aws && go generate

tidy-modules-%:
	@# tidy command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "tidy-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. build-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst tidy-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go mod tidy"

add-module-license-files:
	cd internal/repotools/cmd/eachmodule && \
    	go run . -skip-root \
            "cp $(LICENSE_FILE) ."


################
# Unit Testing #
################

# TODO replace the command with the unit-modules-. once protocol tests pass

unit: lint unit-modules-aws unit-modules-config unit-modules-credentials unit-modules-ec2imds unit-modules-service 
unit-race: lint unit-race-modules-aws unit-race-modules-config unit-race-modules-credentials unit-race-modules-ec2imds unit-race-modules-service

unit-test: test-modules-aws test-modules-config test-modules-credentials test-modules-ec2imds test-modules-service
unit-race-test: test-race-modules-aws test-race-modules-config test-race-modules-credentials test-race-modules-ec2imds test-race-modules-service

unit-race-modules-%:
	@# unit command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "unit-race-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. unit-race-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst unit-race-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go vet ${BUILD_TAGS} --all ./..." \
		"go test ${BUILD_TAGS} ${RUN_NONE} ./..." \
		"go test -timeout=1m ${UNIT_TEST_TAGS} -race -cpu=4 ./..."


unit-modules-%:
	@# unit command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "unit-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. unit-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst unit-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go vet ${BUILD_TAGS} --all ./..." \
		"go test ${BUILD_TAGS} ${RUN_NONE} ./..." \
		"go test -timeout=1m ${UNIT_TEST_TAGS} ./..."

build: build-modules-.

build-modules-%:
	@# build command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "build-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. build-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst build-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go test ${BUILD_TAGS} ${RUN_NONE} ./..."

test: test-modules-.

test-race-modules-%:
	@# Test command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "test-race-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. test-race-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst test-race-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go test -timeout=1m ${UNIT_TEST_TAGS} -race -cpu=4 ./..."

test-modules-%:
	@# Test command that uses the pattern to define the root path that the
	@# module testing will start from. Strips off the "test-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. test-modules-internal_protocoltest
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst test-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go test -timeout=1m ${UNIT_TEST_TAGS} ./..."

##############
# CI Testing #
##############
ci-test: generate unit-race ci-test-generate-validate
ci-test-no-generate: unit-race

ci-test-generate-validate:
	@echo "CI test validate no generated code changes"
	git update-index --assume-unchanged go.mod go.sum
	git add . -A
	gitstatus=`git diff --cached --ignore-space-change`; \
	echo "$$gitstatus"; \
	if [ "$$gitstatus" != "" ] && [ "$$gitstatus" != "skipping validation" ]; then echo "$$gitstatus"; exit 1; fi
	git update-index --no-assume-unchanged go.mod go.sum

#######################
# Integration Testing #
#######################
integration: integ-modules-service

integ-modules-%:
	@# integration command that uses the pattern to define the root path that
	@# the module testing will start from. Strips off the "integ-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. test-modules-service_dynamodb
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst integ-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go test -timeout=10m -tags "integration" -v ${RUN_INTEG} -count 1 ./..."

cleanup-integ-buckets:
	@echo "Cleaning up SDK integraiton resources"
	go run -tags "integration" ./internal/awstesting/cmd/bucket_cleanup/main.go "aws-sdk-go-integration"

##############
# Benchmarks #
##############
bench: bench-modules-.

bench-modules-%:
	@# benchmark command that uses the pattern to define the root path that
	@# the module testing will start from. Strips off the "bench-modules-" and
	@# replaces all "_" with "/".
	@#
	@# e.g. bench-modules-service_dynamodb
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst bench-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go test -timeout=10m -bench . --benchmem ${BUILD_TAGS} ${RUN_NONE} ./..."

##################
# Linting/Verify #
##################
verify: lint vet sdkv1check

lint:
	@echo "go lint SDK and vendor packages"
	@lint=`golint ./...`; \
	dolint=`echo "$$lint" | grep -E -v \
	-e ${LINTIGNORESINGLEFIGHT}`; \
	echo "$$dolint"; \
	if [ "$$dolint" != "" ]; then exit 1; fi

vet: vet-modules-.

vet-modules-%:
	cd ./internal/repotools/cmd/eachmodule \
		&& go run . -p $(subst _,/,$(subst vet-modules-,,$@)) ${EACHMODULE_CONCURRENCY_FLAG} ${EACHMODULE_FAILFAST_FLAG} \
		"go vet ${BUILD_TAGS} --all ./..."

sdkv1check:
	@echo "Checking for usage of AWS SDK for Go v1"
	@sdkv1usage=`go list -test -f '''{{ if not .Standard }}{{ range $$_, $$name := .Imports }} * {{ $$.ImportPath }} -> {{ $$name }}{{ print "\n" }}{{ end }}{{ range $$_, $$name := .TestImports }} *: {{ $$.ImportPath }} -> {{ $$name }}{{ print "\n" }}{{ end }}{{ end}}''' ./... | sort -u | grep '''/aws-sdk-go/'''`; \
	echo "$$sdkv1usage"; \
	if [ "$$sdkv1usage" != "" ]; then exit 1; fi

###################
# Sandbox Testing #
###################
sandbox-tests: sandbox-test-go1.15 sandbox-test-gotip

sandbox-build-go1.15:
	docker build -f ./internal/awstesting/sandbox/Dockerfile.test.go1.15 -t "aws-sdk-go-v2-1.15" .
sandbox-go1.15: sandbox-build-go1.15
	docker run -i -t aws-sdk-go-v2-1.15 bash
sandbox-test-go1.15: sandbox-build-go1.15
	docker run -t aws-sdk-go-v2-1.15

sandbox-build-gotip:
	@echo "Run make update-aws-golang-tip, if this test fails because missing aws-golang:tip container"
	docker build -f ./internal/awstesting/sandbox/Dockerfile.test.gotip -t "aws-sdk-go-v2-tip" .
sandbox-gotip: sandbox-build-gotip
	docker run -i -t aws-sdk-go-v2-tip bash
sandbox-test-gotip: sandbox-build-gotip
	docker run -t aws-sdk-go-v2-tip

update-aws-golang-tip:
	docker build --no-cache=true -f ./internal/awstesting/sandbox/Dockerfile.golang-tip -t "aws-golang:tip" .
