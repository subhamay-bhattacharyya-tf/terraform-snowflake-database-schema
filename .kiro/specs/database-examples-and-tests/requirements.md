# Requirements Document

## Introduction

This document specifies the requirements for populating the Terraform Snowflake database-schema module examples and adding corresponding Terratest tests. The goal is to provide working example configurations that demonstrate various usage patterns of the `database-schema` submodule and ensure test coverage for database and schema creation scenarios.

## Glossary

- **Database_Schema_Module**: The Terraform submodule located at `modules/database-schema/` that creates Snowflake databases and schemas using the `database_configs` variable
- **Example_Configuration**: A complete, working Terraform configuration in the `examples/` directory that demonstrates a specific usage pattern of the module
- **Terratest**: A Go testing framework used to validate Terraform infrastructure code by applying configurations and verifying resources in Snowflake
- **Database_Config**: A map entry in the `database_configs` variable that defines a Snowflake database and its optional schemas

## Requirements

### Requirement 1: Database-Only Example

**User Story:** As a Terraform user, I want an example that creates a single database without schemas, so that I can understand the minimal configuration needed.

#### Acceptance Criteria

1. WHEN a user applies the database-only example, THE Example_Configuration SHALL create exactly one Snowflake database with no schemas
2. THE Example_Configuration SHALL expose the database name and fully qualified name as outputs
3. THE Example_Configuration SHALL use the same provider configuration pattern as the existing basic example
4. THE Example_Configuration SHALL include a descriptive comment in the main.tf file explaining the example purpose

### Requirement 2: Database with One Schema Example

**User Story:** As a Terraform user, I want an example that creates a database with a single schema, so that I can understand how to configure schemas within a database.

#### Acceptance Criteria

1. WHEN a user applies the database-with-one-schema example, THE Example_Configuration SHALL create one database containing exactly one schema
2. THE Example_Configuration SHALL demonstrate optional schema properties including comment and is_managed settings
3. THE Example_Configuration SHALL expose both database and schema names as outputs
4. THE Example_Configuration SHALL expose schema fully qualified names as outputs

### Requirement 3: Database with Multiple Schemas Example

**User Story:** As a Terraform user, I want an example that creates a single database with multiple schemas, so that I can understand how to organize multiple schemas within one database.

#### Acceptance Criteria

1. WHEN a user applies the databases-with-multiple-schemas example, THE Example_Configuration SHALL create one database containing at least three schemas
2. THE Example_Configuration SHALL demonstrate different schema configurations including transient and managed access schemas
3. THE Example_Configuration SHALL expose a nested map of schema names organized by database key
4. THE Example_Configuration SHALL demonstrate the data_retention_time_in_days property on at least one schema

### Requirement 4: Multiple Databases with Multiple Schemas Example

**User Story:** As a Terraform user, I want an example that creates multiple databases each with multiple schemas, so that I can understand complex multi-database configurations.

#### Acceptance Criteria

1. WHEN a user applies the multiple-databases-with-multiple-schemas example, THE Example_Configuration SHALL create at least two databases each containing at least two schemas
2. THE Example_Configuration SHALL demonstrate different database configurations including transient databases
3. THE Example_Configuration SHALL expose all database and schema outputs in a structured format
4. THE Example_Configuration SHALL demonstrate the full range of configurable properties across databases and schemas

### Requirement 5: Single Database Test

**User Story:** As a module maintainer, I want a Terratest that validates single database creation, so that I can ensure the module works correctly for basic use cases.

#### Acceptance Criteria

1. WHEN the single database test runs, THE Terratest SHALL apply the database-only example with a unique database name
2. THE Terratest SHALL verify the database exists in Snowflake after apply
3. THE Terratest SHALL verify the database properties match the configuration
4. THE Terratest SHALL clean up the database after the test completes using terraform destroy
5. THE Terratest SHALL use the same authentication pattern as existing warehouse tests

### Requirement 6: Database with Schema Test

**User Story:** As a module maintainer, I want a Terratest that validates database and schema creation together, so that I can ensure schema creation works correctly.

#### Acceptance Criteria

1. WHEN the database with schema test runs, THE Terratest SHALL apply the database-with-one-schema example with unique names
2. THE Terratest SHALL verify both the database and schema exist in Snowflake after apply
3. THE Terratest SHALL verify schema properties including the parent database relationship
4. THE Terratest SHALL clean up all resources after the test completes

### Requirement 7: Multiple Databases Test

**User Story:** As a module maintainer, I want a Terratest that validates multiple database creation with schemas, so that I can ensure the module handles complex configurations correctly.

#### Acceptance Criteria

1. WHEN the multiple databases test runs, THE Terratest SHALL apply the multiple-databases-with-multiple-schemas example with unique names
2. THE Terratest SHALL verify all databases exist in Snowflake after apply
3. THE Terratest SHALL verify all schemas exist and are associated with the correct databases
4. THE Terratest SHALL verify at least one transient database is created correctly
5. THE Terratest SHALL clean up all resources after the test completes

### Requirement 8: Test Helper Functions

**User Story:** As a test developer, I want helper functions for database and schema verification, so that I can write consistent and maintainable tests.

#### Acceptance Criteria

1. THE Terratest helpers SHALL include a function to check if a database exists in Snowflake
2. THE Terratest helpers SHALL include a function to check if a schema exists in a specific database
3. THE Terratest helpers SHALL include a function to fetch database properties for verification
4. THE Terratest helpers SHALL include a function to fetch schema properties for verification
5. THE Terratest helpers SHALL follow the same patterns as existing warehouse helper functions
