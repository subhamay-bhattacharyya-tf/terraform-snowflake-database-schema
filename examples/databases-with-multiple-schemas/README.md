# Database with Multiple Schemas Example

This example demonstrates how to create a single Snowflake database with multiple schemas using the `database-schema` module. It showcases different schema configurations including transient schemas, managed access, and custom data retention settings.

## Usage

```hcl
module "database" {
  source = "../../modules/database-schema"

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
