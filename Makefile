DEFAULT_BRANCH := main

define _build
	GOOS=$1 GOARCH=amd64 go build -o envcontainer cmd/envcontainer/*.go
endef

.PHONY: compact/linux
compact/linux:
	$(call _build,linux)
	@zip envcontainer_linux_amd64.zip envcontainer 

.PHONY: build
build:
	$(call _build,linux)

.PHONY: run
run:
	go run cmd/envcontainer/*.go

.PHONY: bump-version/major
bump-version/major:  ## Increment the major version (X.y.z)
	bump2version major

.PHONY: bump-version/minor
bump-version/minor:  ## Increment the minor version (x.Y.z)
	bump2version minor

.PHONY:  bump-version/patch
bump-version/patch:  ## Increment the patch version (x.y.Z)
	bump2version patch

.PHONY: release
release:  ## Push the new project version
	git push --follow-tags origin $(DEFAULT_BRANCH)