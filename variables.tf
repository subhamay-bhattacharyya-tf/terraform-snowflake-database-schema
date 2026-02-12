# -----------------------------------------------------------------------------
# Terraform Snowflake Database Schema Module - Variables
# -----------------------------------------------------------------------------
# Input variables for configuring Snowflake databases and schemas.
# -----------------------------------------------------------------------------

variable "database_configs" {
  description = "Map of configuration objects for Snowflake databases and their schemas"
  type = map(object({
    name                        = string
    comment                     = optional(string, null)
    data_retention_time_in_days = optional(number, 1)
    is_transient                = optional(bool, false)
    schemas = optional(list(object({
      name                        = string
      comment                     = optional(string, null)
      is_transient                = optional(bool, false)
      is_managed                  = optional(bool, false)
      data_retention_time_in_days = optional(number, null)
    })), [])
  }))
  default = {}

  validation {
    condition = alltrue([
      for db in var.database_configs : length(db.name) > 0
    ])
    error_message = "Database name must not be empty."
  }

  validation {
    condition = alltrue([
      for db in var.database_configs : alltrue([
        for schema in db.schemas : length(schema.name) > 0
      ])
    ])
    error_message = "Schema name must not be empty."
  }

  validation {
    condition = alltrue([
      for db in var.database_configs : db.data_retention_time_in_days >= 0
    ])
    error_message = "data_retention_time_in_days must be >= 0."
  }

  validation {
    condition = alltrue([
      for db in var.database_configs : alltrue([
        for schema in db.schemas : coalesce(schema.data_retention_time_in_days, 0) >= 0
      ])
    ])
    error_message = "Schema data_retention_time_in_days must be >= 0 or null."
  }
}
