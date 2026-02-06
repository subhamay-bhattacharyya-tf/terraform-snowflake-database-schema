# Changelog

All notable changes to this project will be documented in this file.

## [1.0.1](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-warehouse/compare/v1.0.0...v1.0.1) (2026-02-06)

### Bug Fixes

* **snowflake:** update JWT authenticator to SNOWFLAKE_JWT and remove extra blank line ([1069d20](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-warehouse/commit/1069d20cb9aa25f30eb69d770b3e7fc406194f44))

## [unreleased]

### 🐛 Bug Fixes

- *(snowflake)* Update JWT authenticator to SNOWFLAKE_JWT and remove extra blank line

### 🚜 Refactor

- *(test)* Migrate to gosnowflake config builder for JWT authentication
- *(test)* Improve warehouse property fetching and remove extra blank line

### 📚 Documentation

- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]

### 🎨 Styling

- *(snowflake-warehouse)* Add blank line after module header comment

### ⚙️ Miscellaneous Tasks

- *(testing)* Migrate from Jest to Terratest and restructure examples
- *(github-actions)* Migrate authentication to key-pair and remove property tests
- *(github-actions)* Enhance Terratest output visibility and remove conditional gate
- *(test)* Update Go dependencies and add go.sum
- *(testing)* Migrate to key-pair authentication and add go mod tidy
- *(github-actions)* Add pipefail option to Terratest commands and update Snowflake provider source
## [1.0.0] - 2026-02-04

### 🚀 Features

- *(snowflake-warehouse)* Support multiple warehouses via map configuration

### 🚜 Refactor

- Restructure project to modular Terraform architecture

### 📚 Documentation

- *(readme)* Update badges to reflect Snowflake focus
- Update CHANGELOG.md [skip ci]

### 🎨 Styling

- *(snowflake-warehouse)* Add periods to output descriptions

### ⚙️ Miscellaneous Tasks

- *(github-actions)* Add permissions and token for changelog generation
- *(release)* Version 1.0.0 [skip ci]
