// +build !cgo

package sqlite

import "log"
import "database/sql"

// Stub for builds with no SQLite.
func Init() *sql.DB {
	log.Panic("SQLite not supported in this build.")
	return nil
}