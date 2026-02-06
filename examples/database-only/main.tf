# Example: Single Snowflake Database (No Schemas)
#
# This example demonstrates how to use the database-schema module
# to create a single Snowflake database without any schemas.
# This is the minimal configuration needed to create a database.

module "database" {
  source = "../../modules/database-schema"

  database_configs = var.database_configs
}
