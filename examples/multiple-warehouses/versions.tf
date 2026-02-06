terraform {
  required_version = ">= 1.3.0"

  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = ">= 0.87.0"
    }
  }
}

# Provider configuration using key-pair authentication
# Required environment variables:
#   SNOWFLAKE_ORGANIZATION_NAME - Snowflake organization name
#   SNOWFLAKE_ACCOUNT_NAME      - Snowflake account name
#   SNOWFLAKE_USER              - Snowflake username
#   SNOWFLAKE_ROLE              - Snowflake role
#   SNOWFLAKE_PRIVATE_KEY       - Snowflake private key (PEM format)

provider "snowflake" {
  organization_name = var.snowflake_organization_name
  account_name      = var.snowflake_account_name
  user              = var.snowflake_user
  role              = var.snowflake_role
  authenticator     = "SNOWFLAKE_JWT"
  private_key       = var.snowflake_private_key
}
