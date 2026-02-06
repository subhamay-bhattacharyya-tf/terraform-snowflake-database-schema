// File: test/single_database_test.go
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

// TestSingleDatabase tests creating a single database without schemas
// Property 1: Database Creation Round-Trip
// Property 3: Configuration Fidelity
func TestSingleDatabase(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)

	tfDir := "../examples/database-only"

	databaseConfigs := map[string]interface{}{
		"test_db": map[string]interface{}{
			"name":                        dbName,
			"comment":                     "Terratest single database test",
			"data_retention_time_in_days": 1,
			"is_transient":                false,
			"schemas":                     []interface{}{},
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
	exists := databaseExists(t, db, dbName)
	require.True(t, exists, "Expected database %q to exist in Snowflake", dbName)

	// Property 3: Configuration Fidelity
	props := fetchDatabaseProps(t, db, dbName)
	require.Equal(t, dbName, props.Name)
	require.Contains(t, props.Comment, "Terratest single database test")
}
