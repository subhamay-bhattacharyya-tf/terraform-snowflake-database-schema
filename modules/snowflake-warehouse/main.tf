# Snowflake Warehouse Resource
# Creates and manages one or more Snowflake warehouses based on the warehouse_configs map

resource "snowflake_warehouse" "this" {
  for_each = var.warehouse_configs

  name                      = each.value.name
  warehouse_size            = each.value.warehouse_size
  warehouse_type            = each.value.warehouse_type
  auto_resume               = each.value.auto_resume
  auto_suspend              = each.value.auto_suspend
  initially_suspended       = each.value.initially_suspended
  min_cluster_count         = each.value.min_cluster_count
  max_cluster_count         = each.value.max_cluster_count
  scaling_policy            = each.value.scaling_policy
  enable_query_acceleration = each.value.enable_query_acceleration
  comment                   = each.value.comment
}
