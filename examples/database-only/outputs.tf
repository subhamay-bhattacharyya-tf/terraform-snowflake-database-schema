output "database_names" {
  description = "Map of database config keys to database names"
  value       = module.database.database_names
}

output "database_fully_qualified_names" {
  description = "Map of database config keys to fully qualified names"
  value       = module.database.database_fully_qualified_names
}
