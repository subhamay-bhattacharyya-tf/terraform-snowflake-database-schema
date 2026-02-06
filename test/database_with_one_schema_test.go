// File: test/database_with_schema_test.go
package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestDatabaseWithSchema tests creating a database with one schema
// Property 1: Database Creation Round-Trip
// Property 2: Schema Creation Round-Trip
// Property 3: Configuration Fidelity
func TestDatabaseWithSchema(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := fmt.Sprintf("TT_SCHEMA_%s", unique)

	tfDir := "../examples/database-with-one-schema"

	databaseConfigs := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    dbName,
			"comment": "Terratest database with schema test",
			"schemas": []interface{}{
				map[string]interface{}{
					"name":       schemaName,
					"comment":    "Terratest schema",
					"is_managed": true,
				},
			},
		},
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"database_configs":            databaseConfigs,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer terraform.Destroy(t, tfOptions)
	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	db := openSnowflake(t)
	defer func() { _ = db.Close() }()

	// Property 1: Database Creation Round-Trip
	dbExists := databaseExists(t, db, dbName)
	require.True(t, dbExists, "Expected database %q to exist in Snowflake", dbName)

	// Property 2: Schema Creation Round-Trip
	schemaExistsResult := schemaExists(t, db, dbName, schemaName)
	require.True(t, schemaExistsResult, "Expected schema %q to exist in database %q", schemaName, dbName)

	// Property 3: Configuration Fidelity
	dbProps := fetchDatabaseProps(t, db, dbName)
	require.Equal(t, dbName, dbProps.Name)
	require.Contains(t, dbProps.Comment, "Terratest database with schema test")

	schemaProps := fetchSchemaProps(t, db, dbName, schemaName)
	require.Equal(t, schemaName, schemaProps.Name)
	require.Equal(t, dbName, schemaProps.DatabaseName)
	require.Contains(t, schemaProps.Comment, "Terratest schema")
	require.True(t, schemaProps.IsManagedAccess, "Expected schema to have managed access enabled")
}
