# Database Only Example

This example demonstrates how to create a single Snowflake database without any schemas using the `database-schema` module.

## Usage

```hcl
module "database" {
  source = "../../modules/database-schema"

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

## Running the Example

```bash
terraform init
terraform plan
terraform apply
```
