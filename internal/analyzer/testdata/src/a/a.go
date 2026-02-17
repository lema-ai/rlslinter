package a

import (
	"gorm.io/gorm"
)

// Test case 1: GORM .Row() should be flagged
func testRow(db *gorm.DB) {
	db.Row() // want "GORM .Row\\(\\) is not supported"
}

// Test case 2: GORM .Rows() should be flagged
func testRows(db *gorm.DB) {
	db.Rows() // want "GORM .Rows\\(\\) is not supported"
}

// Test case 3: GORM .Scan() should be flagged
func testScan(db *gorm.DB) {
	var result string
	db.Scan(&result) // want "GORM .Scan\\(\\) is not supported"
}

// Test case 4: Chained GORM .Row() should be flagged
func testChainedRow(db *gorm.DB) {
	db.Where("id = ?", 1).Row() // want "GORM .Row\\(\\) is not supported"
}

// Test case 5: Suppression with nolint:rlslinter should NOT be flagged
func testSuppressionSpecific(db *gorm.DB) {
	//nolint:rlslinter
	db.Row()
}

// Test case 6: Suppression with generic nolint should NOT be flagged
func testSuppressionGeneric(db *gorm.DB) {
	//nolint
	db.Rows()
}

// Test case 7: Inline suppression should NOT be flagged
func testSuppressionInline(db *gorm.DB) {
	db.Scan(&struct{}{}) //nolint:rlslinter
}

// ExcelFile simulates excelize.File to test false positives
type ExcelFile struct{}

func (f *ExcelFile) Rows(sheet string) {}

// Test case 8: Non-GORM .Rows() should NOT be flagged (false positive prevention)
func testExcelRows() {
	file := &ExcelFile{}
	file.Rows("Sheet1") // Should NOT be flagged - not a GORM type
}

// Platform simulates sql.Scanner implementation
type Platform string

// Test case 9: sql.Scanner Scan method implementation should NOT be flagged
func (p *Platform) Scan(value interface{}) error {
	// This is a method declaration, not a call - should NOT be flagged
	return nil
}

// Test case 10: Non-GORM .Scan() should NOT be flagged
func testNonGormScan() {
	var p Platform
	_ = p.Scan("test") // Should NOT be flagged - calling Platform.Scan, not gorm.DB.Scan
}
