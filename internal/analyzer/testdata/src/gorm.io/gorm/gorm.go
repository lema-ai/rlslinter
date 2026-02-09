package gorm

// DB is a minimal stub of gorm.DB for testing
type DB struct{}

// Row is a stub method
func (db *DB) Row() *DB {
	return db
}

// Rows is a stub method
func (db *DB) Rows() *DB {
	return db
}

// Scan is a stub method
func (db *DB) Scan(dest interface{}) *DB {
	return db
}

// Where is a stub method for chaining
func (db *DB) Where(query interface{}, args ...interface{}) *DB {
	return db
}
