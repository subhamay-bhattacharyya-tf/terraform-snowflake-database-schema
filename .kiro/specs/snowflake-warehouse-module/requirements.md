# Requirements Document

## Introduction

This document defines the requirements for a Terraform module that creates and manages Snowflake warehouses. The module accepts a single object variable containing all warehouse configuration properties and creates a Snowflake warehouse resource using the snowflake provider. The module follows the existing repository patterns for Terraform module development.

## Glossary

- **Warehouse_Module**: The Terraform module that creates and manages Snowflake warehouse resources
- **Warehouse_Config**: The input object variable containing all warehouse configuration properties
- **Snowflake_Provider**: The Terraform provider for Snowflake resources (Snowflake-Labs/snowflake)
- **Warehouse_Resource**: The snowflake_warehouse resource created by the module

## Requirements

### Requirement 1: Input Variable Definition

**User Story:** As a Terraform user, I want to provide warehouse configuration as a single object variable, so that I can manage all warehouse settings in one place.

#### Acceptance Criteria

1. THE Warehouse_Module SHALL accept a single input variable named "warehouse_config" of type object
2. WHEN the warehouse_config object is provided, THE Warehouse_Module SHALL support the following properties:
   - name (string, required): The warehouse identifier
   - warehouse_size (string, optional): Size of the warehouse (X-SMALL, SMALL, MEDIUM, LARGE, X-LARGE, 2X-LARGE, 3X-LARGE, 4X-LARGE, 5X-LARGE, 6X-LARGE)
   - warehouse_type (string, optional): Type of warehouse (STANDARD, SNOWPARK-OPTIMIZED)
   - auto_resume (bool, optional): Whether to auto-resume when queries are submitted
   - auto_suspend (number, optional): Seconds of inactivity before auto-suspend
   - initially_suspended (bool, optional): Whether warehouse starts in suspended state
   - min_cluster_count (number, optional): Minimum number of clusters
   - max_cluster_count (number, optional): Maximum number of clusters
   - scaling_policy (string, optional): Scaling policy (STANDARD, ECONOMY)
   - enable_query_acceleration (bool, optional): Whether to enable query acceleration
   - comment (string, optional): Description of the warehouse
3. WHEN optional properties are not provided, THE Warehouse_Module SHALL use Snowflake default values

### Requirement 2: Warehouse Resource Creation

**User Story:** As a Terraform user, I want the module to create a Snowflake warehouse with my specified configuration, so that I can provision warehouses consistently.

#### Acceptance Criteria

1. WHEN a valid warehouse_config is provided, THE Warehouse_Module SHALL create a snowflake_warehouse resource
2. THE Warehouse_Resource SHALL use the name property from warehouse_config as the warehouse identifier
3. WHEN warehouse_size is provided, THE Warehouse_Resource SHALL be created with the specified size
4. WHEN warehouse_type is provided, THE Warehouse_Resource SHALL be created with the specified type
5. WHEN auto_resume is provided, THE Warehouse_Resource SHALL be configured with the specified auto-resume setting
6. WHEN auto_suspend is provided, THE Warehouse_Resource SHALL be configured with the specified auto-suspend timeout
7. WHEN initially_suspended is provided, THE Warehouse_Resource SHALL start in the specified state
8. WHEN min_cluster_count and max_cluster_count are provided, THE Warehouse_Resource SHALL be configured with the specified cluster limits
9. WHEN scaling_policy is provided, THE Warehouse_Resource SHALL use the specified scaling policy
10. WHEN enable_query_acceleration is provided, THE Warehouse_Resource SHALL be configured with the specified query acceleration setting
11. WHEN comment is provided, THE Warehouse_Resource SHALL include the specified comment

### Requirement 3: Input Validation

**User Story:** As a Terraform user, I want the module to validate my input configuration, so that I can catch configuration errors before applying.

#### Acceptance Criteria

1. WHEN warehouse_config.name is empty or null, THE Warehouse_Module SHALL fail with a descriptive error
2. WHEN warehouse_size is provided with an invalid value, THE Warehouse_Module SHALL fail with a descriptive error listing valid sizes
3. WHEN warehouse_type is provided with an invalid value, THE Warehouse_Module SHALL fail with a descriptive error listing valid types
4. WHEN scaling_policy is provided with an invalid value, THE Warehouse_Module SHALL fail with a descriptive error listing valid policies
5. WHEN auto_suspend is provided with a value less than 0, THE Warehouse_Module SHALL fail with a descriptive error
6. WHEN min_cluster_count exceeds max_cluster_count, THE Warehouse_Module SHALL fail with a descriptive error

### Requirement 4: Module Outputs

**User Story:** As a Terraform user, I want the module to output warehouse details, so that I can reference them in other resources or modules.

#### Acceptance Criteria

1. THE Warehouse_Module SHALL output the warehouse name
2. THE Warehouse_Module SHALL output the fully qualified warehouse name
3. THE Warehouse_Module SHALL output the warehouse size
4. THE Warehouse_Module SHALL output the warehouse state (started/suspended)

### Requirement 5: Provider Configuration

**User Story:** As a Terraform user, I want the module to use the Snowflake provider correctly, so that it integrates with my existing Snowflake infrastructure.

#### Acceptance Criteria

1. THE Warehouse_Module SHALL declare the Snowflake-Labs/snowflake provider as a required provider
2. THE Warehouse_Module SHALL specify a compatible provider version constraint
3. THE Warehouse_Module SHALL NOT configure provider authentication (left to the calling module)
