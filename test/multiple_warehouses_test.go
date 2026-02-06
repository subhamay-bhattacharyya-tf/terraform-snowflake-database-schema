// File: test/multiple_warehouses_test.go
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

// TestMultipleWarehouses tests creating multiple warehouses via the module
func TestMultipleWarehouses(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	wh1Name := fmt.Sprintf("TT_ADHOC_%s", unique)
	wh2Name := fmt.Sprintf("TT_LOAD_%s", unique)
	wh3Name := fmt.Sprintf("TT_TRANSFORM_%s", unique)

	tfDir := "../examples/multiple-warehouses"

	warehouseConfigs := map[string]interface{}{
		"adhoc_wh": map[string]interface{}{
			"name":                      wh1Name,
			"warehouse_size":            "X-SMALL",
			"warehouse_type":            "STANDARD",
			"auto_resume":               true,
			"auto_suspend":              60,
			"initially_suspended":       true,
			"min_cluster_count":         1,
			"max_cluster_count":         1,
			"scaling_policy":            "STANDARD",
			"enable_query_acceleration": false,
			"comment":                   "Terratest adhoc warehouse",
		},
		"load_wh": map[string]interface{}{
			"name":                      wh2Name,
			"warehouse_size":            "X-SMALL",
			"warehouse_type":            "STANDARD",
			"auto_resume":               true,
			"auto_suspend":              60,
			"initially_suspended":       true,
			"min_cluster_count":         1,
			"max_cluster_count":         1,
			"scaling_policy":            "STANDARD",
			"enable_query_acceleration": false,
			"comment":                   "Terratest load warehouse",
		},
		"transform_wh": map[string]interface{}{
			"name":                      wh3Name,
			"warehouse_size":            "SMALL",
			"warehouse_type":            "STANDARD",
			"auto_resume":               true,
			"auto_suspend":              120,
			"initially_suspended":       true,
			"min_cluster_count":         1,
			"max_cluster_count":         2,
			"scaling_policy":            "STANDARD",
			"enable_query_acceleration": true,
			"comment":                   "Terratest transform warehouse",
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

	// Verify all three warehouses exist
	for _, whName := range []string{wh1Name, wh2Name, wh3Name} {
		exists := warehouseExists(t, db, whName)
		require.True(t, exists, "Expected warehouse %q to exist in Snowflake", whName)
	}

	// Verify properties of transform warehouse (has different settings)
	props := fetchWarehouseProps(t, db, wh3Name)
	require.Equal(t, wh3Name, props.Name)
	require.Equal(t, "Small", props.Size)
	require.Contains(t, props.Comment, "Terratest transform warehouse")
}
