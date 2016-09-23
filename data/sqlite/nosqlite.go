// +build !cgo

package sqlite

import (
	"database/sql"
	"fmt"
)

// Stub for builds with no SQLite.
func Init() (*sql.DB, error) {
	return nil, fmt.Errorf("SQLite not supported in this build.")
}