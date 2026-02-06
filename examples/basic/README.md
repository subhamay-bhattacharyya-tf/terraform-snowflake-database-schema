# Basic Example - Single Warehouse

This example demonstrates how to create a single Snowflake warehouse using the module.

## Usage

```hcl
module "warehouse" {
  source = "../../modules/snowflake-warehouse"

  warehouse_configs = {
    "my_warehouse" = {
      name                      = "MY_WAREHOUSE"
      warehouse_size            = "X-SMALL"
      warehouse_type            = "STANDARD"
      auto_resume               = true
      auto_suspend              = 60
      initially_suspended       = true
      min_cluster_count         = 1
      max_cluster_count         = 1
      scaling_policy            = "STANDARD"
      enable_query_acceleration = false
      comment                   = "My test warehouse"
    }
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| warehouse_configs | Map of warehouse configuration objects | map(object) | yes |

## Outputs

| Name | Description |
|------|-------------|
| warehouse_names | The names of the created warehouses |
| warehouse_fully_qualified_names | The fully qualified names of the warehouses |
| warehouse_sizes | The sizes of the warehouses |
| warehouse_states | The states of the warehouses |
