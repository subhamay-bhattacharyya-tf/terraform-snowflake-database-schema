variable "database_configs" {
  description = "Map of configuration objects for Snowflake databases and their schemas"
  type = map(object({
    name                        = string
    comment                     = optional(string, null)
    data_retention_time_in_days = optional(number, 1)
    is_transient                = optional(bool, false)
    grants = optional(object({
      usage_roles = optional(list(string), [])
    }), {})
    schemas = optional(list(object({
      name                        = string
      comment                     = optional(string, null)
      is_transient                = optional(bool, false)
      is_managed                  = optional(bool, false)
      data_retention_time_in_days = optional(number, null)
      grants = optional(object({
        usage_roles                    = optional(list(string), [])
        create_file_format_roles       = optional(list(string), [])
        create_stage_roles             = optional(list(string), [])
        create_table_roles             = optional(list(string), [])
        create_pipe_roles              = optional(list(string), [])
        create_dynamic_table_roles     = optional(list(string), [])
        create_stream_roles            = optional(list(string), [])
        create_task_roles              = optional(list(string), [])
        create_view_roles              = optional(list(string), [])
        create_materialized_view_roles = optional(list(string), [])
      }), {})
    })), [])
  }))
  default = {
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

# Snowflake authentication variables
variable "snowflake_organization_name" {
  description = "Snowflake organization name"
  type        = string
  default     = null
}

variable "snowflake_account_name" {
  description = "Snowflake account name"
  type        = string
  default     = null
}

variable "snowflake_user" {
  description = "Snowflake username"
  type        = string
  default     = null
}

variable "snowflake_role" {
  description = "Snowflake role"
  type        = string
  default     = null
}

variable "snowflake_private_key" {
  description = "Snowflake private key for key-pair authentication"
  type        = string
  sensitive   = true
  default     = null
}
