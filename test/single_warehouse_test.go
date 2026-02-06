// File: test/single_warehouse_test.go
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

// TestSingleWarehouse tests creating a single warehouse via the module
func TestSingleWarehouse(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	warehouseName := fmt.Sprintf("TT_WH_%s", unique)

	tfDir := "../examples/basic"

	warehouseConfigs := map[string]interface{}{
		"test_wh": map[string]interface{}{
			"name":                      warehouseName,
			"warehouse_size":            "X-SMALL",
			"warehouse_type":            "STANDARD",
			"auto_resume":               true,
			"auto_suspend":              60,
			"initially_suspended":       true,
			"min_cluster_count":         1,
			"max_cluster_count":         1,
			"scaling_policy":            "STANDARD",
			"enable_query_acceleration": false,
			"comment":                   "Terratest single warehouse test",
		},
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"warehouse_configs":           warehouseConfigs,
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

	exists := warehouseExists(t, db, warehouseName)
	require.True(t, exists, "Expected warehouse %q to exist in Snowflake", warehouseName)

	props := fetchWarehouseProps(t, db, warehouseName)
	require.Equal(t, warehouseName, props.Name)
	require.Equal(t, "X-Small", props.Size)
	require.Contains(t, props.Comment, "Terratest single warehouse test")
}
