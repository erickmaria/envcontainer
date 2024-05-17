define _build
	GOOS=$1 GOARCH=amd64 go build -o envcontainer cmd/envcontainer/envcontainer.go
endef

.PHONY: compact/linux
compact/linux:
	$(call _build,linux)
	@zip envcontainer_v1.0.0_linux_amd64.zip envcontainer 