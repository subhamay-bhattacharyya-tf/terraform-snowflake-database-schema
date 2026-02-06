// File: test/helpers_test.go
package test

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

type WarehouseProps struct {
	Name    string
	Size    string
	Comment string
}

type DatabaseProps struct {
	Name                    string
	Comment                 string
	DataRetentionTimeInDays int
	IsTransient             bool
}

type SchemaProps struct {
	Name                    string
	DatabaseName            string
	Comment                 string
	IsTransient             bool
	IsManagedAccess         bool
	DataRetentionTimeInDays int
}

func openSnowflake(t *testing.T) *sql.DB {
	t.Helper()

	orgName := mustEnv(t, "SNOWFLAKE_ORGANIZATION_NAME")
	accountName := mustEnv(t, "SNOWFLAKE_ACCOUNT_NAME")
	user := mustEnv(t, "SNOWFLAKE_USER")
	privateKeyPEM := mustEnv(t, "SNOWFLAKE_PRIVATE_KEY")
	role := os.Getenv("SNOWFLAKE_ROLE")

	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	require.NotNil(t, block, "Failed to decode PEM block from private key")

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS8 first, then PKCS1
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		require.NoError(t, err, "Failed to parse private key")
	} else {
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		require.True(t, ok, "Private key is not RSA")
	}

	// Build account identifier: orgname-accountname
	account := fmt.Sprintf("%s-%s", orgName, accountName)

	config := gosnowflake.Config{
		Account:       account,
		User:          user,
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}

	if role != "" {
		config.Role = role
	}

	dsn, err := gosnowflake.DSN(&config)
	require.NoError(t, err, "Failed to build DSN")

	db, err := sql.Open("snowflake", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func warehouseExists(t *testing.T, db *sql.DB, warehouseName string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW WAREHOUSES LIKE '%s';", escapeLike(warehouseName))
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func fetchWarehouseProps(t *testing.T, db *sql.DB, warehouseName string) WarehouseProps {
	t.Helper()

	// SHOW WAREHOUSES returns columns in a specific order. We need to scan all columns
	// to get name (col 0), size (col 3), and comment (col 10 in newer versions).
	// Using a simpler approach: query rows and scan by column name using rows.Columns()
	q := fmt.Sprintf("SHOW WAREHOUSES LIKE '%s';", escapeLike(warehouseName))
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	require.NoError(t, err)

	// Find column indices for name, size, comment
	nameIdx, sizeIdx, commentIdx := -1, -1, -1
	for i, col := range cols {
		switch col {
		case "name":
			nameIdx = i
		case "size":
			sizeIdx = i
		case "comment":
			commentIdx = i
		}
	}
	require.NotEqual(t, -1, nameIdx, "name column not found in SHOW WAREHOUSES output")
	require.NotEqual(t, -1, sizeIdx, "size column not found in SHOW WAREHOUSES output")
	require.NotEqual(t, -1, commentIdx, "comment column not found in SHOW WAREHOUSES output")

	require.True(t, rows.Next(), "No warehouse found matching %s", warehouseName)

	// Create slice to hold all column values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	require.NoError(t, err)

	// Extract the values we need
	getName := func(v interface{}) string {
		if v == nil {
			return ""
		}
		if s, ok := v.(string); ok {
			return s
		}
		if b, ok := v.([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", v)
	}

	return WarehouseProps{
		Name:    getName(values[nameIdx]),
		Size:    getName(values[sizeIdx]),
		Comment: getName(values[commentIdx]),
	}
}

func mustEnv(t *testing.T, key string) string {
	t.Helper()
	v := strings.TrimSpace(os.Getenv(key))
	require.NotEmpty(t, v, "Missing required environment variable %s", key)
	return v
}

func escapeLike(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func databaseExists(t *testing.T, db *sql.DB, databaseName string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW DATABASES LIKE '%s';", escapeLike(databaseName))
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func schemaExists(t *testing.T, db *sql.DB, databaseName, schemaName string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW SCHEMAS LIKE '%s' IN DATABASE %s;", escapeLike(schemaName), databaseName)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func fetchDatabaseProps(t *testing.T, db *sql.DB, databaseName string) DatabaseProps {
	t.Helper()

	q := fmt.Sprintf("SHOW DATABASES LIKE '%s';", escapeLike(databaseName))
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	require.NoError(t, err)

	nameIdx, commentIdx, retentionIdx, transientIdx := -1, -1, -1, -1
	for i, col := range cols {
		switch col {
		case "name":
			nameIdx = i
		case "comment":
			commentIdx = i
		case "retention_time":
			retentionIdx = i
		case "is_transient":
			transientIdx = i
		}
	}
	require.NotEqual(t, -1, nameIdx, "name column not found")

	require.True(t, rows.Next(), "No database found matching %s", databaseName)

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	require.NoError(t, err)

	props := DatabaseProps{
		Name: getString(values[nameIdx]),
	}
	if commentIdx != -1 {
		props.Comment = getString(values[commentIdx])
	}
	if retentionIdx != -1 {
		props.DataRetentionTimeInDays = getInt(values[retentionIdx])
	}
	if transientIdx != -1 {
		props.IsTransient = getString(values[transientIdx]) == "true"
	}

	return props
}

func fetchSchemaProps(t *testing.T, db *sql.DB, databaseName, schemaName string) SchemaProps {
	t.Helper()

	q := fmt.Sprintf("SHOW SCHEMAS LIKE '%s' IN DATABASE %s;", escapeLike(schemaName), databaseName)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	require.NoError(t, err)

	nameIdx, dbIdx, commentIdx, retentionIdx, transientIdx, managedIdx := -1, -1, -1, -1, -1, -1
	for i, col := range cols {
		switch col {
		case "name":
			nameIdx = i
		case "database_name":
			dbIdx = i
		case "comment":
			commentIdx = i
		case "retention_time":
			retentionIdx = i
		case "is_transient":
			transientIdx = i
		case "is_managed_access":
			managedIdx = i
		}
	}
	require.NotEqual(t, -1, nameIdx, "name column not found")

	require.True(t, rows.Next(), "No schema found matching %s in database %s", schemaName, databaseName)

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	require.NoError(t, err)

	props := SchemaProps{
		Name: getString(values[nameIdx]),
	}
	if dbIdx != -1 {
		props.DatabaseName = getString(values[dbIdx])
	}
	if commentIdx != -1 {
		props.Comment = getString(values[commentIdx])
	}
	if retentionIdx != -1 {
		props.DataRetentionTimeInDays = getInt(values[retentionIdx])
	}
	if transientIdx != -1 {
		props.IsTransient = getString(values[transientIdx]) == "true"
	}
	if managedIdx != -1 {
		props.IsManagedAccess = getString(values[managedIdx]) == "true"
	}

	return props
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

func getInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	case []byte:
		var i int
		fmt.Sscanf(string(val), "%d", &i)
		return i
	}
	return 0
}
