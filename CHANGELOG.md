# Changelog
All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [v1.x.x] 2025-06-18
### Added 
- Added `sitehost_stack_name` resource.
- Added `sitehost_stack` resource.
- Added `sitehost_stack_environment` resource.
- Added `sitehost_cloud_database` resource.
- Added `sitehost_cloud_database_user` resource.
- Added `sitehost_cloud_database_grant` resource.
- Added `sitehost_cloud_ssh_user` resource.

- Added `sitehost_cloud_database` data source.
- Added `sitehost_cloud_databases` data source.
- Added `sitehost_cloud_database_grant` data source.
- Added `sitehost_cloud_ssh_user` data source.
- Added `sitehost_server` data source.
- Added `sitehost_stack` data source.
- Added `sitehost_stacks` data source.
- Added `sitehost_stack_environment` data source.
- Added `sitehost_ssh_key` data source.
- Added `sitehost_ssh_keys` data source.

### Fixed

### Updated
- Update to GoSH v0.6.0

## [v1.3.0] 2025-06-12
### Added
- Added `sitehost_server_firewall` resource.
- Added `sitehost_server_security_group` resource.

### Updated
- Updated GoSH version to v0.5.0.

## [v1.2.0] 2024-03-13
### Added
- Added `sitehost_ssh_key` resource.

### Updated
- Updated GoSH version to v0.3.4.

## [v1.1.0] 2023-03-27
### Added
- Added `sitehost_dns_zone` resource.
- Added `sitehost_dns_record` resource.
- Added DNS example.
### Updated
- Updated GoSH version to v0.3.2.

## [v1.0.1] 2022-12-09
### Added
- GitHub PR actions to run go vet, go lint and go mod tidy.

### Updated
- Updated GoSH version.
- Updated our code to conform with Golang linter.
- Updated our README to link to our Golang style and our license.
- Upgraded golang v1.19.3.
- Upgraded golangci-lint v1.50.1.

## [v1.0.0] 2022-07-15
## Added
- Initial release SiteHost Terraform provider. Support Job and Server API endpoint.
