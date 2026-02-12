# -----------------------------------------------------------------------------
# Terraform Snowflake Database Schema Module - Outputs
# -----------------------------------------------------------------------------
# Output values for created Snowflake databases and schemas.
# -----------------------------------------------------------------------------

output "database_names" {
  description = "Map of database config keys to database names."
  value       = { for k, v in snowflake_database.this : k => v.name }
}

output "database_fully_qualified_names" {
  description = "Map of database config keys to fully qualified names."
  value       = { for k, v in snowflake_database.this : k => v.fully_qualified_name }
}

output "databases" {
  description = "All database resource objects."
  value       = snowflake_database.this
}

output "schema_names" {
  description = "Nested map of database keys to schema names to schema name values."
  value = {
    for db_key in distinct([for k, v in local.schemas : v.db_key]) : db_key => {
      for k, v in local.schemas : v.schema.name => snowflake_schema.this[k].name
      if v.db_key == db_key
    }
  }
}

output "schema_fully_qualified_names" {
  description = "Nested map of database keys to schema names to fully qualified names."
  value = {
    for db_key in distinct([for k, v in local.schemas : v.db_key]) : db_key => {
      for k, v in local.schemas : v.schema.name => snowflake_schema.this[k].fully_qualified_name
      if v.db_key == db_key
    }
  }
}

output "schemas" {
  description = "All schema resource objects."
  value       = snowflake_schema.this
}
