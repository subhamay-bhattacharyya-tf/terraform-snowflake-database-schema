variable "warehouse_configs" {
  description = "Map of configuration objects for Snowflake warehouses"
  type = map(object({
    name                      = string
    warehouse_size            = optional(string, "X-SMALL")
    warehouse_type            = optional(string, "STANDARD")
    auto_resume               = optional(bool, true)
    auto_suspend              = optional(number, 60)
    initially_suspended       = optional(bool, true)
    min_cluster_count         = optional(number, 1)
    max_cluster_count         = optional(number, 1)
    scaling_policy            = optional(string, "STANDARD")
    enable_query_acceleration = optional(bool, false)
    comment                   = optional(string, null)
  }))
  default = {}

  validation {
    condition     = alltrue([for k, wh in var.warehouse_configs : length(wh.name) > 0])
    error_message = "Warehouse name must not be empty."
  }

  validation {
    condition = alltrue([for k, wh in var.warehouse_configs : contains([
      "X-SMALL", "XSMALL", "SMALL", "MEDIUM", "LARGE"
    ], upper(wh.warehouse_size))])
    error_message = "Invalid warehouse_size. Valid values: X-SMALL, SMALL, MEDIUM, LARGE."
  }

  validation {
    condition     = alltrue([for k, wh in var.warehouse_configs : contains(["STANDARD"], upper(wh.warehouse_type))])
    error_message = "Invalid warehouse_type. Valid values: STANDARD."
  }

  validation {
    condition     = alltrue([for k, wh in var.warehouse_configs : contains(["STANDARD", "ECONOMY"], upper(wh.scaling_policy))])
    error_message = "Invalid scaling_policy. Valid values: STANDARD, ECONOMY."
  }

  validation {
    condition     = alltrue([for k, wh in var.warehouse_configs : wh.auto_suspend >= 0])
    error_message = "auto_suspend must be >= 0 seconds."
  }

  validation {
    condition     = alltrue([for k, wh in var.warehouse_configs : wh.min_cluster_count <= wh.max_cluster_count])
    error_message = "min_cluster_count must not exceed max_cluster_count."
  }
}
