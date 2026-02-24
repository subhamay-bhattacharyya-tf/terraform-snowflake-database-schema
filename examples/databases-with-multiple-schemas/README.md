# Database with Multiple Schemas Example

This example demonstrates how to create a single Snowflake database with multiple schemas using the `database-schema` module. It showcases different schema configurations including transient schemas, managed access, custom data retention settings, and grants.

## Usage

```hcl
module "database" {
  source = "../.."

  database_configs = {
    datawarehouse = {
      name                        = "DATA_WAREHOUSE"
      comment                     = "Central data warehouse"
      data_retention_time_in_days = 7
      grants = {
        usage_roles = ["DATA_READER_ROLE", "ETL_ROLE"]
      }
      schemas = [
        {
          name       = "RAW"
          comment    = "Raw ingested data"
          is_managed = false
          grants = {
            usage_roles        = ["ETL_ROLE"]
            create_table_roles = ["ETL_ROLE"]
            create_stage_roles = ["ETL_ROLE"]
            create_pipe_roles  = ["ETL_ROLE"]
          }
        },
        {
          name         = "STAGING"
          comment      = "Data transformation staging area"
          is_transient = true
          grants = {
            usage_roles              = ["ETL_ROLE"]
            create_table_roles       = ["ETL_ROLE"]
            create_dynamic_table_roles = ["ETL_ROLE"]
            create_stream_roles      = ["ETL_ROLE"]
            create_task_roles        = ["ETL_ROLE"]
          }
        },
        {
          name                        = "CURATED"
          comment                     = "Curated business data"
          is_managed                  = true
          data_retention_time_in_days = 14
          grants = {
            usage_roles                    = ["DATA_READER_ROLE", "ETL_ROLE"]
            create_table_roles             = ["ETL_ROLE"]
            create_view_roles              = ["ETL_ROLE"]
            create_materialized_view_roles = ["ETL_ROLE"]
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
| schema_names | Nested map of database keys to schema names |
| schema_fully_qualified_names | Nested map of database keys to schema fully qualified names |

## Running the Example

```bash
terraform init
terraform plan
terraform apply
```
