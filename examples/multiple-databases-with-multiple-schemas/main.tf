# Example: Multiple Snowflake Databases with Multiple Schemas
#
# This example demonstrates how to use the database-schema module
# to create multiple Snowflake databases, each with multiple schemas.
# It showcases the full range of configurable properties including
# transient databases and schemas, managed access, and data retention.

module "database" {
  source = "../.."

  database_configs = var.database_configs
}
