# Terraform Snowflake Module - Database Schema

![Release](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema/actions/workflows/ci.yaml/badge.svg)&nbsp;![Snowflake](https://img.shields.io/badge/Snowflake-29B5E8?logo=snowflake&logoColor=white)&nbsp;![Commit Activity](https://img.shields.io/github/commit-activity/t/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Last Commit](https://img.shields.io/github/last-commit/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Release Date](https://img.shields.io/github/release-date/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Repo Size](https://img.shields.io/github/repo-size/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![File Count](https://img.shields.io/github/directory-file-count/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Issues](https://img.shields.io/github/issues/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Top Language](https://img.shields.io/github/languages/top/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema)&nbsp;![Custom Endpoint](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bsubhamay/4770674f89a7a0961ce5f0bbce0cda1d/raw/terraform-snowflake-database-schema.json?)

A Terraform module for creating and managing Snowflake databases and schemas using a map of configuration objects. Supports creating single or multiple databases with nested schemas in a single module call.

## Features

- Map-based configuration for creating single or multiple databases
- Nested schema configuration within each database
- Built-in input validation with descriptive error messages
- Sensible defaults for optional properties
- Outputs keyed by database identifier for easy reference
- Support for transient databases and schemas
- Support for managed access schemas
- Configurable data retention time at database and schema level

## Usage

### Single Database (No Schemas)

```hcl
module "database" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema"

  database_configs = {
    analytics = {
      name                        = "ANALYTICS_DB"
      comment                     = "Analytics database for reporting"
      data_retention_time_in_days = 1
      is_transient                = false
      schemas                     = []
    }
  }
}
```

### Database with One Schema

```hcl
module "database" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema"

  database_configs = {
    app = {
      name    = "APPLICATION_DB"
      comment = "Main application database"
      schemas = [
        {
          name       = "PUBLIC_DATA"
          comment    = "Public facing data schema"
          is_managed = true
        }
      ]
    }
  }
}
```

### Database with Multiple Schemas

```hcl
module "database" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema"

  database_configs = {
    datawarehouse = {
      name                        = "DATA_WAREHOUSE"
      comment                     = "Central data warehouse"
      data_retention_time_in_days = 7
      schemas = [
        {
          name       = "RAW"
          comment    = "Raw ingested data"
          is_managed = false
        },
        {
          name         = "STAGING"
          comment      = "Data transformation staging area"
          is_transient = true
        },
        {
          name                        = "CURATED"
          comment                     = "Curated business data"
          is_managed                  = true
          data_retention_time_in_days = 14
        }
      ]
    }
  }
}
```

### Multiple Databases with Multiple Schemas

```hcl
module "database" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-database-schema"

  database_configs = {
    production = {
      name    = "PROD_DB"
      comment = "Production database"
      schemas = [
        { name = "APP", comment = "Application schema" },
        { name = "AUDIT", comment = "Audit logging schema", is_managed = true }
      ]
    },
    development = {
      name         = "DEV_DB"
      comment      = "Development database"
      is_transient = true
      schemas = [
        { name = "SANDBOX", comment = "Developer sandbox" },
        { name = "TESTING", comment = "Test data schema", is_transient = true }
      ]
    }
  }
}
```

## Examples

- [Database Only](examples/database-only) - Create a single database without schemas
- [Database with One Schema](examples/database-with-one-schema) - Create a database with a single schema
- [Database with Multiple Schemas](examples/databases-with-multiple-schemas) - Create a database with multiple schemas
- [Multiple Databases with Multiple Schemas](examples/multiple-databases-with-multiple-schemas) - Create multiple databases with multiple schemas

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 0.87.0 |

## Providers

| Name | Version |
|------|---------|
| snowflake | >= 0.87.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| database_configs | Map of configuration objects for Snowflake databases and their schemas | `map(object)` | `{}` | no |

### database_configs Object Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| name | string | - | Database name (required) |
| comment | string | null | Description of the database |
| data_retention_time_in_days | number | 1 | Time Travel data retention period in days |
| is_transient | bool | false | Whether the database is transient |
| schemas | list(object) | [] | List of schema configurations |

### schemas Object Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| name | string | - | Schema name (required) |
| comment | string | null | Description of the schema |
| is_transient | bool | false | Whether the schema is transient |
| is_managed | bool | false | Whether the schema has managed access |
| data_retention_time_in_days | number | null | Time Travel data retention (inherits from database if null) |

## Outputs

| Name | Description |
|------|-------------|
| database_names | Map of database config keys to database names |
| database_fully_qualified_names | Map of database config keys to fully qualified names |
| databases | All database resource objects |
| schema_names | Nested map of database keys to schema names |
| schema_fully_qualified_names | Nested map of database keys to schema fully qualified names |
| schemas | All schema resource objects |

## Validation

The module validates inputs and provides descriptive error messages for:

- Empty database name
- Empty schema name
- Negative data_retention_time_in_days value

## Testing

The module includes Terratest-based integration tests:

```bash
cd test
go mod tidy
go test -v -timeout 30m
```

Required environment variables for testing:
- `SNOWFLAKE_ORGANIZATION_NAME` - Snowflake organization name
- `SNOWFLAKE_ACCOUNT_NAME` - Snowflake account name
- `SNOWFLAKE_USER` - Snowflake username
- `SNOWFLAKE_ROLE` - Snowflake role (e.g., "SYSADMIN")
- `SNOWFLAKE_PRIVATE_KEY` - Snowflake private key for key-pair authentication

### Test Coverage

| Test File | Example Tested | Properties Validated |
|-----------|----------------|---------------------|
| `single_database_test.go` | database-only | Database creation, configuration fidelity |
| `database_with_schema_test.go` | database-with-one-schema | Database/schema creation, managed access |
| `database_with_multiple_schemas_test.go` | databases-with-multiple-schemas | Multiple schemas, transient schema, managed access |
| `multiple_databases_test.go` | multiple-databases-with-multiple-schemas | Multiple databases, transient resources |

## CI/CD Configuration

The CI workflow runs on:
- Push to `main`, `feature/**`, and `bug/**` branches (when `*.tf`, `examples/**`, or `test/**` changes)
- Pull requests to `main` (when `*.tf`, `examples/**`, or `test/**` changes)
- Manual workflow dispatch

The workflow includes:
- Terraform validation and format checking
- Examples validation
- Terratest integration tests (output displayed in GitHub Step Summary)
- Changelog generation (non-main branches)
- Semantic release (main branch only)

The CI workflow uses the following GitHub organization variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `TERRAFORM_VERSION` | Terraform version for CI jobs | `1.3.0` |
| `GO_VERSION` | Go version for Terratest | `1.21` |
| `SNOWFLAKE_ORGANIZATION_NAME` | Snowflake organization name | - |
| `SNOWFLAKE_ACCOUNT_NAME` | Snowflake account name | - |
| `SNOWFLAKE_USER` | Snowflake username | - |
| `SNOWFLAKE_ROLE` | Snowflake role (e.g., SYSADMIN) | - |

The following GitHub secrets are required for Terratest integration tests:

| Secret | Description | Required |
|--------|-------------|----------|
| `SNOWFLAKE_PRIVATE_KEY` | Snowflake private key for key-pair authentication | Yes |

## License

MIT License - See [LICENSE](LICENSE) for details.
