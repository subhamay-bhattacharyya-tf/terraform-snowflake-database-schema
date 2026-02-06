# Example: Snowflake Database with One Schema
#
# This example demonstrates how to use the database-schema module
# to create a Snowflake database with a single schema.
# The schema includes optional properties like comment and managed access.

module "database" {
  source = "../../modules/database-schema"

  database_configs = var.database_configs
}
