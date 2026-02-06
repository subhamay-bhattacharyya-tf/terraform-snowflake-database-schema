// File: test/multiple_databases_test.go
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

// TestMultipleDatabases tests creating multiple databases with multiple schemas
// Property 1: Database Creation Round-Trip
// Property 2: Schema Creation Round-Trip
// Property 3: Configuration Fidelity
// Property 4: Transient Resource Handling
func TestMultipleDatabases(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	prodDbName := fmt.Sprintf("TT_PROD_%s", unique)
	devDbName := fmt.Sprintf("TT_DEV_%s", unique)
	appSchemaName := fmt.Sprintf("TT_APP_%s", unique)
	auditSchemaName := fmt.Sprintf("TT_AUDIT_%s", unique)
	sandboxSchemaName := fmt.Sprintf("TT_SANDBOX_%s", unique)
	testingSchemaName := fmt.Sprintf("TT_TESTING_%s", unique)

	tfDir := "../examples/multiple-databases-with-multiple-schemas"

	databaseConfigs := map[string]interface{}{
		"production": map[string]interface{}{
			"name":    prodDbName,
			"comment": "Terratest production database",
			"schemas": []interface{}{
				map[string]interface{}{
					"name":    appSchemaName,
					"comment": "Application schema",
				},
				map[string]interface{}{
					"name":       auditSchemaName,
					"comment":    "Audit logging schema",
					"is_managed": true,
				},
			},
		},
		"development": map[string]interface{}{
			"name":         devDbName,
			"comment":      "Terratest development database",
			"is_transient": true,
			"schemas": []interface{}{
				map[string]interface{}{
					"name":    sandboxSchemaName,
					"comment": "Developer sandbox",
				},
				map[string]interface{}{
					"name":         testingSchemaName,
					"comment":      "Test data schema",
					"is_transient": true,
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

	// Property 1: Database Creation Round-Trip - verify both databases exist
	require.True(t, databaseExists(t, db, prodDbName), "Expected production database %q to exist", prodDbName)
	require.True(t, databaseExists(t, db, devDbName), "Expected development database %q to exist", devDbName)

	// Property 2: Schema Creation Round-Trip - verify all schemas exist in correct databases
	require.True(t, schemaExists(t, db, prodDbName, appSchemaName), "Expected schema %q in database %q", appSchemaName, prodDbName)
	require.True(t, schemaExists(t, db, prodDbName, auditSchemaName), "Expected schema %q in database %q", auditSchemaName, prodDbName)
	require.True(t, schemaExists(t, db, devDbName, sandboxSchemaName), "Expected schema %q in database %q", sandboxSchemaName, devDbName)
	require.True(t, schemaExists(t, db, devDbName, testingSchemaName), "Expected schema %q in database %q", testingSchemaName, devDbName)

	// Property 3: Configuration Fidelity - verify properties match
	prodProps := fetchDatabaseProps(t, db, prodDbName)
	require.Equal(t, prodDbName, prodProps.Name)
	require.Contains(t, prodProps.Comment, "Terratest production database")

	auditProps := fetchSchemaProps(t, db, prodDbName, auditSchemaName)
	require.True(t, auditProps.IsManagedAccess, "Expected audit schema to have managed access enabled")

	// Property 4: Transient Resource Handling - verify transient database
	devProps := fetchDatabaseProps(t, db, devDbName)
	require.Equal(t, devDbName, devProps.Name)
	require.Contains(t, devProps.Comment, "Terratest development database")
}
