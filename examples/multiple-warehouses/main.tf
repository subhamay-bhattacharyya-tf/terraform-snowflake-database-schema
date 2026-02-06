# Example: Multiple Snowflake Warehouses
#
# This example demonstrates how to use the snowflake-warehouse module
# to create multiple Snowflake warehouses using a map of configurations.

module "warehouses" {
  source = "../../modules/snowflake-warehouse"

  warehouse_configs = var.warehouse_configs
}
