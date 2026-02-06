output "database_names" {
  description = "Map of database config keys to database names"
  value       = module.database.database_names
}

output "database_fully_qualified_names" {
  description = "Map of database config keys to fully qualified names"
  value       = module.database.database_fully_qualified_names
}

output "databases" {
  description = "All database resource objects"
  value       = module.database.databases
}

output "schema_names" {
  description = "Nested map of database keys to schema names"
  value       = module.database.schema_names
}

output "schema_fully_qualified_names" {
  description = "Nested map of database keys to schema fully qualified names"
  value       = module.database.schema_fully_qualified_names
}

output "schemas" {
  description = "All schema resource objects"
  value       = module.database.schemas
}
