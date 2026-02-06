output "warehouse_names" {
  description = "The names of the created warehouses"
  value       = module.warehouse.warehouse_names
}

output "warehouse_fully_qualified_names" {
  description = "The fully qualified names of the warehouses"
  value       = module.warehouse.warehouse_fully_qualified_names
}

output "warehouse_sizes" {
  description = "The sizes of the warehouses"
  value       = module.warehouse.warehouse_sizes
}

output "warehouse_states" {
  description = "The states of the warehouses (STARTED or SUSPENDED)"
  value       = module.warehouse.warehouse_states
}
