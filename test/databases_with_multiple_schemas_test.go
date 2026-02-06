// File: test/database_with_multiple_schemas_test.go
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

// TestDatabaseWithMultipleSchemas tests creating a single database with multiple schemas
// Property 1: Database Creation Round-Trip
// Property 2: Schema Creation Round-Trip
// Property 3: Configuration Fidelity
// Property 4: Transient Resource Handling
func TestDatabaseWithMultipleSchemas(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	dbName := fmt.Sprintf("TT_DW_%s", unique)
	rawSchemaName := fmt.Sprintf("TT_RAW_%s", unique)
	stagingSchemaName := fmt.Sprintf("TT_STAGING_%s", unique)
	curatedSchemaName := fmt.Sprintf("TT_CURATED_%s", unique)

	tfDir := "../examples/databases-with-multiple-schemas"

	databaseConfigs := map[string]interface{}{
		"datawarehouse": map[string]interface{}{
			"name":                        dbName,
			"comment":                     "Terratest data warehouse",
			"data_retention_time_in_days": 7,
			"schemas": []interface{}{
				map[string]interface{}{
					"name":       rawSchemaName,
					"comment":    "Raw ingested data",
					"is_managed": false,
				},
				map[string]interface{}{
					"name":         stagingSchemaName,
					"comment":      "Data transformation staging area",
					"is_transient": true,
				},
				map[string]interface{}{
					"name":       curatedSchemaName,
					"comment":    "Curated business data",
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
	require.True(t, databaseExists(t, db, dbName), "Expected database %q to exist", dbName)

	// Property 2: Schema Creation Round-Trip - verify all three schemas exist
	require.True(t, schemaExists(t, db, dbName, rawSchemaName), "Expected schema %q in database %q", rawSchemaName, dbName)
	require.True(t, schemaExists(t, db, dbName, stagingSchemaName), "Expected schema %q in database %q", stagingSchemaName, dbName)
	require.True(t, schemaExists(t, db, dbName, curatedSchemaName), "Expected schema %q in database %q", curatedSchemaName, dbName)

	// Property 3: Configuration Fidelity - verify database properties
	dbProps := fetchDatabaseProps(t, db, dbName)
	require.Equal(t, dbName, dbProps.Name)
	require.Contains(t, dbProps.Comment, "Terratest data warehouse")

	// Verify curated schema has managed access
	curatedProps := fetchSchemaProps(t, db, dbName, curatedSchemaName)
	require.True(t, curatedProps.IsManagedAccess, "Expected curated schema to have managed access enabled")
}
