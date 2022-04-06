.DEFAULT_GOAL := all

DOCKER_REGISTRY           ?= registry.digitalocean.com/ice-io
DOCKER_TAG                ?= latest-locally
GO_VERSION_MANIFEST       := https://raw.githubusercontent.com/actions/go-versions/main/versions-manifest.json
REQUIRED_COVERAGE_PERCENT := 60
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

generate:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -d cmd/freezer -g freezer.go -o cmd/freezer/api
	swag fmt -d cmd/freezer -g freezer.go
#	go install github.com/golang/mock/mockgen@latest
#	mockgen -source=CHANGE_ME.go -destination=CHANGE_ME.go -package=CHANGE_ME

checkGenerated: generate
	@if git status --porcelain | grep -e [.]go -e [.]json -e [.]yaml; then \
		echo "Please commit generated files, using 'make generate'."; \
		git --no-pager diff; \
		exit 1; \
	fi; \
	true;

build:
	set -xe; \
	[ -d cmd ] && find ./cmd -mindepth 1 -maxdepth 1 -type d -print | while read service; do \
		go build -tags=go_json -race -v -o $${service}/bin $${service}; \
	done; true;

buildAllSupportedPlatforms: clean
	echo "Linux"
	for arch in arm64 amd64 s390x ppc64le riscv64; do \
		$(MAKE) CGO_ENABLED=0 GOOS=linux GOARCH=$$arch binary-for-each-service clean; \
	done;
	echo "Darwin"
	for arch in arm64 amd64; do \
		$(MAKE) CGO_ENABLED=0 GOOS=darwin GOARCH=$$arch binary-for-each-service clean; \
	done;
	echo "Windows"
	for arch in arm64 amd64; do \
		$(MAKE) CGO_ENABLED=0 GOOS=windows GOARCH=$$arch binary-for-each-service clean; \
	done;

binary-for-each-service:
	set -xe; \
	[ -d cmd ] && find ./cmd -mindepth 1 -maxdepth 1 -type d -print | while read service; do \
		echo "$@: $${service} / $(GOARCH)" ; \
		go build -tags=go_json -a -v -o $${service}/bin $${service}; \
	done; true;

binary-specific-service:
	set -xe; \
	echo "$@: $(SERVICE_NAME) / $(GOARCH)" ; \
	go build -tags=go_json -a -v -o ./cmd/$${SERVICE_NAME}/bin ./cmd/$${SERVICE_NAME}; \

test:
	@go version
	# TODO make -race work
	#go test -tags=go_json -v -race -cover -coverprofile=$(COVERAGE_FILE) -covermode atomic ./...
	go test -tags=go_json -v -cover -coverprofile=$(COVERAGE_FILE) -covermode atomic ./...
	@grep -v "_generated.go" $(COVERAGE_FILE) > tmp$(COVERAGE_FILE)
	@mv -f tmp$(COVERAGE_FILE) $(COVERAGE_FILE)

# TODO should be improved to a per file check and maybe against a previous value
#(maybe we should use something like SonarQube for this?)
coverage: $(COVERAGE_FILE)
	@t=`go tool cover -func=$(COVERAGE_FILE) | grep total | grep -Eo '[0-9]+\.[0-9]+'`;\
	echo "Total coverage: $${t}%"; \
	if [ "$${t%.*}" -lt $(REQUIRED_COVERAGE_PERCENT) ]; then \
		echo "ERROR: It has to be at least $(REQUIRED_COVERAGE_PERCENT)%"; \
		exit 1; \
	fi;

benchmark:
	# TODO make -race work
	go test -tags=go_json -run=^$ -v -bench=. -benchmem -benchtime 10s ./cmd/freezer
	go test -tags=go_json -run=^$ -v -bench=. -benchmem -benchtime 10s ./cmd/freezer-refrigerant

clean:
	@go clean
	@rm -f tmp$(COVERAGE_FILE) $(COVERAGE_FILE) 2>/dev/null || true
	@test -d cmd && find ./cmd -mindepth 2 -maxdepth 2 -type f -name bin -exec rm -f {} \; || true;

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
	find ./cmd -mindepth 1 -maxdepth 1 -type d -print | while read service; do \
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
	go run -v local.go

getAddLicense:
	GO111MODULE=off go get -v -u github.com/google/addlicense

addLicense: getAddLicense
	`go env GOPATH`/bin/addlicense -f LICENSE.header * .github/* .deploy/*

checkLicense: getAddLicense
	`go env GOPATH`/bin/addlicense -f LICENSE.header -check * .github/* .deploy/*

all: checkLicense checkModVersion checkIfAllDependenciesAreUpToDate checkGenerated build buildAllSupportedPlatforms test coverage benchmark clean
local: addLicense checkLicense updateGoModVersion updateAllDependencies generate build buildMultiPlatformDockerImage test coverage benchmark lint clean
dockerfile: binary-specific-service
