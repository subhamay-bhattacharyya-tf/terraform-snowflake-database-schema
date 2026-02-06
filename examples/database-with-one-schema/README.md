# Database with One Schema Example

This example demonstrates how to create a Snowflake database with a single schema using the `database-schema` module. The schema includes managed access configuration.

## Usage

```hcl
module "database" {
  source = "../../modules/database-schema"

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
