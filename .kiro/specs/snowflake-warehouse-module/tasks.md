# Implementation Plan: Snowflake Warehouse Module

## Overview

This plan implements a Terraform module for creating Snowflake warehouses. The module accepts a single object variable with all warehouse configuration properties and creates a `snowflake_warehouse` resource using the Snowflake-Labs/snowflake provider.

## Tasks

- [x] 1. Create module directory structure and provider configuration
  - [x] 1.1 Create `modules/snowflake-warehouse/` directory structure
    - Create versions.tf with Snowflake-Labs/snowflake provider requirement (>= 0.87.0)
    - Specify terraform required_version >= 1.3.0
    - _Requirements: 5.1, 5.2, 5.3_

- [x] 2. Implement input variable with validation
  - [x] 2.1 Create variables.tf with warehouse_config object variable
    - Define typed object with all properties (name, warehouse_size, warehouse_type, auto_resume, auto_suspend, initially_suspended, min_cluster_count, max_cluster_count, scaling_policy, enable_query_acceleration, comment)
    - Set optional properties with default values
    - _Requirements: 1.1, 1.2, 1.3_
  
  - [x] 2.2 Add validation blocks for input constraints
    - Validate name is not empty
    - Validate warehouse_size against valid enum values
    - Validate warehouse_type against valid enum values
    - Validate scaling_policy against valid enum values
    - Validate auto_suspend >= 0
    - Validate min_cluster_count <= max_cluster_count
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

  - [x] 2.3 Write property test for valid input acceptance
    - **Property 3: Valid Input Acceptance**
    - **Validates: Requirements 1.1, 2.1**

  - [x] 2.4 Write property test for invalid input rejection
    - **Property 4: Invalid Input Rejection**
    - **Validates: Requirements 3.1, 3.2, 3.3, 3.4, 3.5, 3.6**

- [x] 3. Implement warehouse resource
  - [x] 3.1 Create main.tf with snowflake_warehouse resource
    - Map all warehouse_config properties to resource arguments
    - Use var.warehouse_config.* for all attribute values
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9, 2.10, 2.11_

  - [x] 3.2 Write property test for configuration passthrough
    - **Property 1: Configuration Passthrough**
    - **Validates: Requirements 1.2, 2.1-2.11**

  - [x] 3.3 Write property test for default value application
    - **Property 2: Default Value Application**
    - **Validates: Requirements 1.3**

- [x] 4. Implement module outputs
  - [x] 4.1 Create outputs.tf with warehouse outputs
    - Output warehouse_name from resource
    - Output warehouse_fully_qualified_name from resource
    - Output warehouse_size from resource
    - Output warehouse_state based on initially_suspended
    - _Requirements: 4.1, 4.2, 4.3, 4.4_

  - [x] 4.2 Write property test for output correctness
    - **Property 5: Output Correctness**
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.4**

- [x] 5. Create example usage
  - [x] 5.1 Create examples/snowflake-warehouse/main.tf with example configuration
    - Show usage with the provided example configuration (SN_TEST_STREAMLIT_WH)
    - Include provider configuration placeholder
    - _Requirements: 1.2, 2.1_

- [x] 6. Checkpoint - Validate module structure
  - Run `terraform init` and `terraform validate` on the module
  - Ensure all tests pass, ask the user if questions arise

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- The module follows existing repository patterns (similar to S3 bucket module structure)
- Property tests use Terratest framework for Terraform testing
- Integration tests require a Snowflake account and are not included in this plan
