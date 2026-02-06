# Example: Single Snowflake Warehouse
#
# This example demonstrates how to use the snowflake-warehouse module
# to create a single Snowflake warehouse.

module "warehouse" {
  source = "../../modules/snowflake-warehouse"

  warehouse_configs = var.warehouse_configs
}
