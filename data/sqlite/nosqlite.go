// +build !linux OR !cgo

package sqlite

// Stub for builds with no SQLite.
func Init() {
}