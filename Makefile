.DEFAULT_GOAL := all

DOCKER_REGISTRY           ?= registry.digitalocean.com/ice-io
DOCKER_TAG                ?= latest-locally
GO_VERSION_MANIFEST       := https://raw.githubusercontent.com/actions/go-versions/main/versions-manifest.json
REQUIRED_COVERAGE_PERCENT := 0
COVERAGE_FILE             := cover.out
REPOSITORY                := $(shell basename `pwd`)

CGO_ENABLED := 1
GOOS         ?=
GOARCH       ?=
SERVICE_NAME ?=
SERVICES    := $(wildcard ./cmd/*)

export CGO_ENABLED GOOS GOARCH SERVICE_NAME

define getLatestGoPatchVersion
	$(shell curl -s $(GO_VERSION_MANIFEST) | jq -r '.[0].version')
endef

define getLatestGoMinorVersion
	$(shell echo $(call getLatestGoPatchVersion) | cut -f1,2 -d'.')
endef

latestGoVersion:
	@echo $(call getLatestGoPatchVersion)

latestGoMinorVersion:
	@echo $(call getLatestGoMinorVersion)

updateGoModVersion:
	go mod edit -go $(call getLatestGoMinorVersion)

checkModVersion: updateGoModVersion
	@if git status --porcelain | grep -q go.mod; then \
		echo "Outdated go version in go.mod. Please update it using 'make updateGoModVersion' and make sure everything works correctly and tests pass then commit the changes."; \
		exit 1; \
	fi; \
	true;

updateAllDependencies:
	go get -t -u ./...
	go mod tidy

checkIfAllDependenciesAreUpToDate: updateAllDependencies
	@if git status --porcelain | grep -q go.sum; then \
		echo "Some dependencies are outdated. Please update all dependencies using 'make updateAllDependencies' and make sure everything works correctly and tests pass then commit the changes."; \
		exit 1; \
	fi; \
	true;

generate-swagger:
	swag init --parseDependency --parseInternal -d ${SERVICE} -g $(shell echo "$${SERVICE##*/}" | sed 's/-/_/g').go -o ${SERVICE}/api;

generate-swaggers:
	go install github.com/swaggo/swag/cmd/swag@latest
	set -xe; \
	[ -d cmd ] && find ./cmd -mindepth 1 -maxdepth 1 -type d -print | grep -v 'fixture' | sed 's/\.\///g' | while read service; do \
		env SERVICE=$${service} $(MAKE) generate-swagger; \
	done;

format-swagger:
	swag fmt -d ${SERVICE} -g $(shell echo "$${SERVICE##*/}" | sed 's/-/_/g').go

format-swaggers:
	set -xe; \
	[ -d cmd ] && find ./cmd -mindepth 1 -maxdepth 1 -type d -print | grep -v 'fixture' | sed 's/\.\///g' | while read service; do \
		env SERVICE=$${service} $(MAKE) format-swagger; \
	done;

generate-mocks:
#	go install github.com/golang/mock/mockgen@latest
#	mockgen -source=CHANGE_ME.go -destination=CHANGE_ME.go -package=CHANGE_ME

generate:
	$(MAKE) generate-swaggers
	$(MAKE) format-swaggers
	$(MAKE) generate-mocks
	$(MAKE) addLicense
	$(MAKE) format-imports

checkGenerated: generate
	@if git status --porcelain | grep -e [.]go -e [.]json -e [.]yaml; then \
		echo "Please commit generated files, using 'make generate'."; \
		git --no-pager diff; \
		exit 1; \
	fi; \
	true;

build-all@ci/cd:
	go build -tags=go_json -a -v -race ./...

build: build-all@ci/cd

binary-specific-service:
	set -xe; \
	echo "$@: $(SERVICE_NAME) / $(GOARCH)" ; \
	go build -tags=go_json -a -v -o ./cmd/$${SERVICE_NAME}/bin ./cmd/$${SERVICE_NAME}; \

test:
	set -xe; \
	mf="$$(pwd)/Makefile"; \
	find . -mindepth 1 -maxdepth 4 -type d -print | grep -v '\./\.' | grep -v '/\.' | sed 's/\.\///g' | while read service; do \
		cd $${service} ; \
		if [[ $$(find . -mindepth 1 -maxdepth 1 -type f -print | grep -E '_test.go' | wc -l | sed "s/ //g") -gt 0 ]]; then \
			make -f $$mf test@ci/cd; \
		fi ; \
		for ((i=0;i<$$(echo "$${service}" | grep -o "/" | wc -l | sed "s/ //g");i++)); do \
          	cd .. ; \
        done; \
        cd .. ; \
	done;

# TODO should be improved to a per file check and maybe against a previous value
#(maybe we should use something like SonarQube for this?)
coverage: $(COVERAGE_FILE)
	@t=`go tool cover -func=$(COVERAGE_FILE) | grep total | grep -Eo '[0-9]+\.[0-9]+'`;\
	echo "Total coverage: $${t}%"; \
	if [ "$${t%.*}" -lt $(REQUIRED_COVERAGE_PERCENT) ]; then \
		echo "ERROR: It has to be at least $(REQUIRED_COVERAGE_PERCENT)%"; \
		exit 1; \
	fi;

test@ci/cd:
	# TODO make -race work
	go test -timeout 20m -tags=go_json,test -v -cover -coverprofile=$(COVERAGE_FILE) -covermode atomic

benchmark@ci/cd:
	# TODO make -race work
	go test -timeout 20m -tags=go_json,test -run=^$ -v -bench=. -benchmem -benchtime 10s

benchmark:
	set -xe; \
	mf="$$(pwd)/Makefile"; \
	find . -mindepth 1 -maxdepth 4 -type d -print | grep -v '\./\.' | grep -v '/\.' | sed 's/\.\///g' | while read service; do \
		cd $${service} ; \
		if [[ $$(find . -mindepth 1 -maxdepth 1 -type f -print | grep -E '_test.go' | wc -l | sed "s/ //g") -gt 0 ]]; then \
			make -f $$mf benchmark@ci/cd; \
		fi ; \
		for ((i=0;i<$$(echo "$${service}" | grep -o "/" | wc -l | sed "s/ //g");i++)); do \
          	cd .. ; \
        done; \
        cd .. ; \
	done;

print-all-packages-with-tests:
	set -xe; \
	find . -mindepth 1 -maxdepth 4 -type d -print | grep -v '\./\.' | grep -v '/\.' | sed 's/\.\///g' | while read service; do \
		cd $${service} ; \
		if [[ $$(find . -mindepth 1 -maxdepth 1 -type f -print | grep -E '_test.go' | wc -l | sed "s/ //g") -gt 0 ]]; then \
			echo "$${service}"; \
		fi ; \
		for ((i=0;i<$$(echo "$${service}" | grep -o "/" | wc -l | sed "s/ //g");i++)); do \
          	cd .. ; \
        done; \
        cd .. ; \
	done;

clean:
	@go clean
	@rm -f tmp$(COVERAGE_FILE) $(COVERAGE_FILE) 2>/dev/null || true
	@test -d cmd && find ./cmd -mindepth 2 -maxdepth 2 -type f -name bin -exec rm -f {} \; || true;
	@test -d cmd && find ./cmd -mindepth 2 -maxdepth 2 -type d -name bins -exec rm -Rf {} \; || true;
	@find . -name ".tmp-*" -exec rm -Rf {} \; || true;
	@find . -mindepth 1 -maxdepth 3 -type f -name $(COVERAGE_FILE) -exec rm -Rf {} \; || true;
	@find . -mindepth 1 -maxdepth 3 -type f -name tmp$(COVERAGE_FILE) -exec rm -Rf {} \; || true;

lint:
	golangci-lint run

# run specific service by its name
run-%:
	go run -tags=go_json -v ./cmd/$*

run:
ifeq ($(words $(SERVICES)),1)
	$(MAKE) $(subst ./cmd/,run-,$(SERVICES))
else
	@echo "Do not know what to run"
	@echo "Targets:"
	@for target in $(subst ./cmd/,run-,$(SERVICES)); do \
		echo "  $${target}"; \
	done; false;
endif

# run specific service by its name
binary-run-%:
	./cmd/$*/bin

binary-run:
ifeq ($(words $(SERVICES)),1)
	$(MAKE) $(subst ./cmd/,binary-run-,$(SERVICES))
else
	@echo "Do not know what to run"
	@echo "Targets:"
	@for target in $(subst ./cmd/,binary-run-,$(SERVICES)); do \
		echo "  $${target}"; \
	done; false;
endif

# note: it requires make-4.3+ to run that
buildMultiPlatformDockerImage:
	set -xe; \
	find ./cmd -mindepth 1 -maxdepth 1 -type d -print | grep -v 'fixture' | while read service; do \
		for arch in amd64 arm64 s390x ppc64le; do \
			docker buildx build \
				--platform linux/$${arch} \
				-f $${service}/Dockerfile \
				--label os=linux \
				--label arch=$${arch} \
				--force-rm \
				--pull -t $(DOCKER_REGISTRY)/$(REPOSITORY)/$${service##*/}:$(DOCKER_TAG) \
				--build-arg SERVICE_NAME=$${service##*/} \
				--build-arg TARGETARCH=$${arch} \
				--build-arg TARGETOS=linux \
				.; \
			done; \
		done;

start-test-environment:
	#go run -race -v local.go
	go run -v local.go --type all

start-test-environment-%:
	#go run -race -v local.go
	go run -v local.go --type $*

getAddLicense:
	GO111MODULE=off go get -v -u github.com/google/addlicense

addLicense: getAddLicense
	`go env GOPATH`/bin/addlicense -f LICENSE.header * .github/*

checkLicense: getAddLicense
	`go env GOPATH`/bin/addlicense -f LICENSE.header -check * .github/*

fix-field-alignment:
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	fieldalignment -fix ./...

format-imports:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/daixiang0/gci@latest
	gci write -s standard -s default -s "prefix(github.com/ice-blockchain)" ./..
	goimports -w -local github.com/ice-blockchain ./..

print-token-%:
	go run -v local.go --generateAuth $*

start-seeding:
	go run -v local.go --startSeeding true

all: checkLicense checkModVersion checkIfAllDependenciesAreUpToDate checkGenerated build test coverage benchmark clean
local: addLicense updateGoModVersion updateAllDependencies generate build buildMultiPlatformDockerImage test coverage benchmark lint clean
dockerfile: binary-specific-service
