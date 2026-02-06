output "warehouse_names" {
  description = "The names of the created warehouses."
  value       = { for k, v in snowflake_warehouse.this : k => v.name }
}

output "warehouse_fully_qualified_names" {
  description = "The fully qualified names of the warehouses."
  value       = { for k, v in snowflake_warehouse.this : k => v.fully_qualified_name }
}

output "warehouse_sizes" {
  description = "The sizes of the warehouses."
  value       = { for k, v in snowflake_warehouse.this : k => v.warehouse_size }
}

output "warehouse_states" {
  description = "The states of the warehouses (STARTED or SUSPENDED)."
  value       = { for k, wh in var.warehouse_configs : k => wh.initially_suspended ? "SUSPENDED" : "STARTED" }
}

output "warehouses" {
  description = "All warehouse resources."
  value       = snowflake_warehouse.this
}
