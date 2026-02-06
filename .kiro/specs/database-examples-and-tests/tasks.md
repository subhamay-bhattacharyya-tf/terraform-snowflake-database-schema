# Implementation Plan: Database Examples and Tests

## Overview

This plan implements four Terraform example configurations for the database-schema module and adds corresponding Terratest tests. The implementation follows existing patterns from the warehouse examples and tests.

## Tasks

- [-] 1. Implement database-only example
  - [x] 1.1 Create main.tf with module invocation
    - Add module block referencing `../../modules/database-schema`
    - Include descriptive comment explaining the example purpose
    - _Requirements: 1.1, 1.4_
  - [x] 1.2 Create variables.tf with database_configs and auth variables
    - Define database_configs variable matching module interface
    - Add Snowflake authentication variables (organization_name, account_name, user, role, private_key)
    - _Requirements: 1.3_
  - [x] 1.3 Create outputs.tf exposing database outputs
    - Output database_names from module
    - Output database_fully_qualified_names from module
    - _Requirements: 1.2_
  - [ ] 1.4 Create versions.tf with provider configuration
    - Set required Terraform version >= 1.3.0
    - Configure Snowflake provider with JWT authentication
    - _Requirements: 1.3_

- [ ] 2. Implement database-with-one-schema example
  - [ ] 2.1 Create main.tf with module invocation
    - Add module block referencing `../../modules/database-schema`
    - Include descriptive comment explaining the example purpose
    - _Requirements: 2.1_
  - [ ] 2.2 Create variables.tf with database_configs including schema
    - Define database_configs with one schema demonstrating comment and is_managed
    - Add Snowflake authentication variables
    - _Requirements: 2.2_
  - [ ] 2.3 Create outputs.tf exposing database and schema outputs
    - Output database_names and database_fully_qualified_names
    - Output schema_names and schema_fully_qualified_names
    - _Requirements: 2.3, 2.4_
  - [ ] 2.4 Create versions.tf with provider configuration
    - _Requirements: 2.1_

- [ ] 3. Implement databases-with-multiple-schemas example
  - [ ] 3.1 Create main.tf with module invocation
    - Add module block referencing `../../modules/database-schema`
    - _Requirements: 3.1_
  - [ ] 3.2 Create variables.tf with database containing multiple schemas
    - Define database_configs with at least 3 schemas
    - Include transient schema, managed access schema, and schema with data_retention_time_in_days
    - _Requirements: 3.1, 3.2, 3.4_
  - [ ] 3.3 Create outputs.tf exposing nested schema map
    - Output database outputs
    - Output schema_names as nested map by database key
    - Output schema_fully_qualified_names
    - _Requirements: 3.3_
  - [ ] 3.4 Create versions.tf with provider configuration
    - _Requirements: 3.1_

- [ ] 4. Implement multiple-databases-with-multiple-schemas example
  - [ ] 4.1 Create main.tf with module invocation
    - Add module block referencing `../../modules/database-schema`
    - _Requirements: 4.1_
  - [ ] 4.2 Create variables.tf with multiple databases and schemas
    - Define at least 2 databases each with at least 2 schemas
    - Include transient database configuration
    - Demonstrate full range of configurable properties
    - _Requirements: 4.1, 4.2, 4.4_
  - [ ] 4.3 Create outputs.tf exposing all outputs
    - Output all database outputs (names, fully_qualified_names, databases)
    - Output all schema outputs (names, fully_qualified_names, schemas)
    - _Requirements: 4.3_
  - [ ] 4.4 Create versions.tf with provider configuration
    - _Requirements: 4.1_

- [ ] 5. Checkpoint - Verify examples
  - Ensure all examples have valid Terraform syntax
  - Run `terraform validate` on each example directory
  - Ask the user if questions arise

- [ ] 6. Add database and schema helper functions to tests
  - [ ] 6.1 Add DatabaseProps and SchemaProps structs to helpers_test.go
    - Define DatabaseProps with Name, Comment, DataRetentionTimeInDays, IsTransient
    - Define SchemaProps with Name, DatabaseName, Comment, IsTransient, IsManagedAccess, DataRetentionTimeInDays
    - _Requirements: 8.3, 8.4_
  - [ ] 6.2 Implement databaseExists function
    - Query Snowflake using `SHOW DATABASES LIKE` pattern
    - Return boolean indicating existence
    - _Requirements: 8.1_
  - [ ] 6.3 Implement schemaExists function
    - Query Snowflake using `SHOW SCHEMAS IN DATABASE` pattern
    - Return boolean indicating existence
    - _Requirements: 8.2_
  - [ ] 6.4 Implement fetchDatabaseProps function
    - Query Snowflake for database properties
    - Parse and return DatabaseProps struct
    - Follow existing fetchWarehouseProps pattern
    - _Requirements: 8.3, 8.5_
  - [ ] 6.5 Implement fetchSchemaProps function
    - Query Snowflake for schema properties
    - Parse and return SchemaProps struct
    - Follow existing fetchWarehouseProps pattern
    - _Requirements: 8.4, 8.5_

- [ ] 7. Implement single database test
  - [ ] 7.1 Create single_database_test.go
    - Generate unique database name with TT_DB_ prefix
    - Configure terraform options for database-only example
    - Apply configuration and verify database exists
    - Verify database properties match configuration
    - **Property 1: Database Creation Round-Trip**
    - **Property 3: Configuration Fidelity**
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ] 8. Implement database with schema test
  - [ ] 8.1 Create database_with_schema_test.go
    - Generate unique names for database and schema
    - Configure terraform options for database-with-one-schema example
    - Apply configuration and verify database and schema exist
    - Verify schema is associated with correct database
    - Verify schema properties (is_managed, comment)
    - **Property 1: Database Creation Round-Trip**
    - **Property 2: Schema Creation Round-Trip**
    - **Property 3: Configuration Fidelity**
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [ ] 9. Implement multiple databases test
  - [ ] 9.1 Create multiple_databases_test.go
    - Generate unique names for all databases and schemas
    - Configure terraform options for multiple-databases-with-multiple-schemas example
    - Apply configuration and verify all databases exist
    - Verify all schemas exist and are associated with correct databases
    - Verify transient database is created correctly
    - **Property 1: Database Creation Round-Trip**
    - **Property 2: Schema Creation Round-Trip**
    - **Property 3: Configuration Fidelity**
    - **Property 4: Transient Resource Handling**
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 10. Final checkpoint - Ensure all tests pass
  - Run `go test -v ./...` in test directory
  - Verify all tests pass against Snowflake
  - Ensure all resources are cleaned up after tests
  - Ask the user if questions arise

## Notes

- All examples follow the same file structure pattern as the existing basic example
- Tests use unique resource names with TT_ prefix to avoid conflicts
- Tests run in parallel using `t.Parallel()`
- All tests use `defer terraform.Destroy()` to ensure cleanup
- Property tests are integrated into the Terratest assertions rather than separate property-based testing library
