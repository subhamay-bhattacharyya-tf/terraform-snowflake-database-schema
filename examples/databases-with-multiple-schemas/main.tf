# Example: Snowflake Database with Multiple Schemas
#
# This example demonstrates how to use the database-schema module
# to create a single Snowflake database with multiple schemas.
# It showcases different schema configurations including transient,
# managed access, and custom data retention settings.

module "database" {
  source = "../.."

  database_configs = var.database_configs
}
