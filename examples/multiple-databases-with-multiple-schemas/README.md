# Multiple Databases with Multiple Schemas Example

This example demonstrates how to create multiple Snowflake databases, each with multiple schemas, using the `database-schema` module. It showcases the full range of configurable properties including transient databases and schemas, managed access, grants, and various schema configurations.

## Usage

```hcl
module "database" {
  source = "../.."

  database_configs = {
    production = {
      name    = "PROD_DB"
      comment = "Production database"
      grants = {
        usage_roles = ["PROD_READER_ROLE", "PROD_WRITER_ROLE"]
      }
      schemas = [
        {
          name    = "APP"
          comment = "Application schema"
          grants = {
            usage_roles                    = ["PROD_READER_ROLE", "PROD_WRITER_ROLE"]
            create_table_roles             = ["PROD_WRITER_ROLE"]
            create_view_roles              = ["PROD_WRITER_ROLE"]
            create_materialized_view_roles = ["PROD_WRITER_ROLE"]
          }
        },
        {
          name       = "AUDIT"
          comment    = "Audit logging schema"
          is_managed = true
          grants = {
            usage_roles        = ["AUDIT_ROLE"]
            create_table_roles = ["AUDIT_ROLE"]
            create_stream_roles = ["AUDIT_ROLE"]
          }
        }
      ]
    },
    development = {
      name         = "DEV_DB"
      comment      = "Development database"
      is_transient = true
      grants = {
        usage_roles = ["DEV_ROLE"]
      }
      schemas = [
        {
          name    = "SANDBOX"
          comment = "Developer sandbox"
          grants = {
            usage_roles                    = ["DEV_ROLE"]
            create_table_roles             = ["DEV_ROLE"]
            create_view_roles              = ["DEV_ROLE"]
            create_dynamic_table_roles     = ["DEV_ROLE"]
            create_stream_roles            = ["DEV_ROLE"]
            create_task_roles              = ["DEV_ROLE"]
          }
        },
        {
          name         = "TESTING"
          comment      = "Test data schema"
          is_transient = true
          grants = {
            usage_roles        = ["DEV_ROLE"]
            create_table_roles = ["DEV_ROLE"]
          }
        }
      ]
    }
  }
}
```

## Supported Schema Grants

| Grant | Description |
|-------|-------------|
| usage_roles | USAGE privilege on the schema |
| create_file_format_roles | CREATE FILE FORMAT privilege |
| create_stage_roles | CREATE STAGE privilege |
| create_table_roles | CREATE TABLE privilege |
| create_pipe_roles | CREATE PIPE privilege |
| create_dynamic_table_roles | CREATE DYNAMIC TABLE privilege |
| create_stream_roles | CREATE STREAM privilege |
| create_task_roles | CREATE TASK privilege |
| create_view_roles | CREATE VIEW privilege |
| create_materialized_view_roles | CREATE MATERIALIZED VIEW privilege |

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 0.87.0 |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| database_configs | Map of database configurations | `map(object)` | yes |
| snowflake_organization_name | Snowflake organization name | `string` | yes |
| snowflake_account_name | Snowflake account name | `string` | yes |
| snowflake_user | Snowflake username | `string` | yes |
| snowflake_role | Snowflake role | `string` | yes |
| snowflake_private_key | Snowflake private key for authentication | `string` | yes |

## Outputs

| Name | Description |
|------|-------------|
| database_names | Map of database config keys to database names |
| database_fully_qualified_names | Map of database config keys to fully qualified names |
| databases | All database resource objects |
| schema_names | Nested map of database keys to schema names |
| schema_fully_qualified_names | Nested map of database keys to schema fully qualified names |
| schemas | All schema resource objects |

## Running the Example

```bash
terraform init
terraform plan
terraform apply
```
