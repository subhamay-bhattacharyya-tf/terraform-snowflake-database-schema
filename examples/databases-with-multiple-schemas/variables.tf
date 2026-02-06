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
  default = {
    datawarehouse = {
      name                        = "DATA_WAREHOUSE"
      comment                     = "Central data warehouse"
      data_retention_time_in_days = 7
      schemas = [
        {
          name       = "RAW"
          comment    = "Raw ingested data"
          is_managed = false
        },
        {
          name         = "STAGING"
          comment      = "Data transformation staging area"
          is_transient = true
        },
        {
          name                        = "CURATED"
          comment                     = "Curated business data"
          is_managed                  = true
          data_retention_time_in_days = 14
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
