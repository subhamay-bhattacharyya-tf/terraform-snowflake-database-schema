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
