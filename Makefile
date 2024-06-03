define _build
	GOOS=$1 GOARCH=amd64 go build -o envcontainer cmd/envcontainer/envcontainer.go
endef

.PHONY: compact/linux
compact/linux:
	$(call _build,linux)
	@zip envcontainer_v1.0.0_linux_amd64.zip envcontainer 


.PHONY: bump-version/major
bump-version/major:  ## Increment the major version (X.y.z)
	bump2version major

.PHONY: bump-version/minor
bump-version/minor:  ## Increment the minor version (x.Y.z)
	bump2version minor

.PHONY:  bump-version/patch
bump-version/patch:  ## Increment the patch version (x.y.Z)
	bump2version patch