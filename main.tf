# -----------------------------------------------------------------------------
# Terraform Snowflake Database Schema Module
# -----------------------------------------------------------------------------
# This module creates Snowflake databases and schemas using a map-based
# configuration. It supports creating single or multiple databases with
# nested schemas in a single module call.
# -----------------------------------------------------------------------------

locals {
  schemas = merge([
    for db_key, db in var.database_configs : {
      for schema in db.schemas :
      "${db_key}.${schema.name}" => {
        db_key        = db_key
        database_name = db.name
        schema        = schema
      }
    }
  ]...)

  # Flatten database grants for iteration
  database_usage_grants = merge([
    for db_key, db in var.database_configs : {
      for role in db.grants.usage_roles :
      "${db_key}_${role}" => {
        db_key = db_key
        role   = role
      }
    }
  ]...)

  # Flatten schema grants for iteration
  schema_usage_grants = merge([
    for schema_key, schema_data in local.schemas : {
      for role in schema_data.schema.grants.usage_roles :
      "${schema_key}_${role}" => {
        schema_key = schema_key
        role       = role
      }
    }
  ]...)

  schema_create_file_format_grants = merge([
    for schema_key, schema_data in local.schemas : {
      for role in schema_data.schema.grants.create_file_format_roles :
      "${schema_key}_${role}" => {
        schema_key = schema_key
        role       = role
      }
    }
  ]...)

  schema_create_stage_grants = merge([
    for schema_key, schema_data in local.schemas : {
      for role in schema_data.schema.grants.create_stage_roles :
      "${schema_key}_${role}" => {
        schema_key = schema_key
        role       = role
      }
    }
  ]...)

  schema_create_table_grants = merge([
    for schema_key, schema_data in local.schemas : {
      for role in schema_data.schema.grants.create_table_roles :
      "${schema_key}_${role}" => {
        schema_key = schema_key
        role       = role
      }
    }
  ]...)

  schema_create_pipe_grants = merge([
    for schema_key, schema_data in local.schemas : {
      for role in schema_data.schema.grants.create_pipe_roles :
      "${schema_key}_${role}" => {
        schema_key = schema_key
        role       = role
      }
    }
  ]...)
}

resource "snowflake_database" "this" {
  for_each = var.database_configs

  name                        = each.value.name
  comment                     = each.value.comment
  data_retention_time_in_days = each.value.data_retention_time_in_days
  is_transient                = each.value.is_transient
}

resource "snowflake_schema" "this" {
  for_each = local.schemas

  name                        = each.value.schema.name
  database                    = snowflake_database.this[each.value.db_key].name
  comment                     = each.value.schema.comment
  is_transient                = each.value.schema.is_transient
  with_managed_access         = each.value.schema.is_managed
  data_retention_time_in_days = each.value.schema.data_retention_time_in_days
}

# -----------------------------------------------------------------------------
# Database Grants
# -----------------------------------------------------------------------------

# Database USAGE grants
resource "snowflake_grant_privileges_to_account_role" "database_usage" {
  for_each = local.database_usage_grants

  privileges        = ["USAGE"]
  account_role_name = each.value.role

  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.this[each.value.db_key].fully_qualified_name
  }
}

# -----------------------------------------------------------------------------
# Schema Grants
# -----------------------------------------------------------------------------

# Schema USAGE grants
resource "snowflake_grant_privileges_to_account_role" "schema_usage" {
  for_each = local.schema_usage_grants

  privileges        = ["USAGE"]
  account_role_name = each.value.role

  on_schema {
    schema_name = snowflake_schema.this[each.value.schema_key].fully_qualified_name
  }
}

# Schema CREATE FILE FORMAT grants
resource "snowflake_grant_privileges_to_account_role" "schema_create_file_format" {
  for_each = local.schema_create_file_format_grants

  privileges        = ["CREATE FILE FORMAT"]
  account_role_name = each.value.role

  on_schema {
    schema_name = snowflake_schema.this[each.value.schema_key].fully_qualified_name
  }
}

# Schema CREATE STAGE grants
resource "snowflake_grant_privileges_to_account_role" "schema_create_stage" {
  for_each = local.schema_create_stage_grants

  privileges        = ["CREATE STAGE"]
  account_role_name = each.value.role

  on_schema {
    schema_name = snowflake_schema.this[each.value.schema_key].fully_qualified_name
  }
}

# Schema CREATE TABLE grants
resource "snowflake_grant_privileges_to_account_role" "schema_create_table" {
  for_each = local.schema_create_table_grants

  privileges        = ["CREATE TABLE"]
  account_role_name = each.value.role

  on_schema {
    schema_name = snowflake_schema.this[each.value.schema_key].fully_qualified_name
  }
}

# Schema CREATE PIPE grants
resource "snowflake_grant_privileges_to_account_role" "schema_create_pipe" {
  for_each = local.schema_create_pipe_grants

  privileges        = ["CREATE PIPE"]
  account_role_name = each.value.role

  on_schema {
    schema_name = snowflake_schema.this[each.value.schema_key].fully_qualified_name
  }
}


output "debug_database_usage_grants" {
  value = local.database_usage_grants
}

output "debug_schema_usage_grants" {
  value = local.schema_usage_grants
}