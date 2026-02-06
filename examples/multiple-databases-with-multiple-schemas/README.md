# Multiple Databases with Multiple Schemas Example

This example demonstrates how to create multiple Snowflake databases, each with multiple schemas, using the `database-schema` module. It showcases the full range of configurable properties including transient databases and schemas, managed access, and various schema configurations.

## Usage

```hcl
module "database" {
  source = "../../modules/database-schema"

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
