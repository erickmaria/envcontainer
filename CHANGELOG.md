## Unreleased

## v2.0.0 - 2024-08-25

### Added

- chore: remove volume when stop container
- chore: remove option user on .envcontainer.yaml
- chore: add new flag to start with vscode
- feat: add new flag to start with vscode
- feat: add commant to list envcontainers 
- fix: regex to validate ports pattern

### Fixed

- Error in github workflow for uploading artifacts when creating release

## v1.1.0 - 2024-06-03

### Added

- [#21](https://github.com/erickmaria/envcontainer/pull/21): Create automation to create new releases and tags
- [#19](https://github.com/erickmaria/envcontainer/pull/19): Improvement container volumes

### Fixed

- [6406042](https://github.com/erickmaria/envcontainer/commit/64060422ea0c5abe6b87bfdfa82f5b1026ffa40b): Error building and finding container image when project name has space

## v1.0.0 - 2024-05-17

### Added

- Release stable version

### Fixed

- [#19](https://github.com/erickmaria/envcontainer/pull/19): Container not restart when stopped
