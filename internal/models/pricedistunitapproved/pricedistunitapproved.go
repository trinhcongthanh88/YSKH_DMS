package pricedistunitapproved

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "YSKH_DMS/database"


	// "github.com/davecgh/go-spew/spew"
)

// =====================
// ENTITY - MODEL PRICEDISTUNITAPPROVED
// =====================

type PriceDistUnitApproved struct {
	PriceEntryId 		   int 		  	   `json:"-"`
	PriListCode            string      	   `json:"priListCode"`
	DistUnitCode           string          `json:"distUnitCode"`
    DistUnitName           string          `json:"distUnitName"`
	FromDate		   	   string 		   `json:"fromDate"`
	ToDate		      	   string 		   `json:"toDate"`
}



// =====================
// HELPER: FORMAT DATE CHO MSSQL
// =====================
func formatDateForMSSQL(dateStr string) (any, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" ||
		strings.EqualFold(dateStr, "null") ||
		strings.EqualFold(dateStr, "nil") ||
		strings.EqualFold(dateStr, "na") ||
		strings.EqualFold(dateStr, "<null>") {
		return nil, nil // Không lỗi, chỉ là không có giá trị
	}
	

	// Possible input layouts, in order of priority
	layouts := []string{
		// ISO with T (e.g., 2025-04-15T13:45:00)
		"2006-01-02T15:04:05",
		// Date only ISO (e.g., 2025-04-15)
		"2006-01-02",
		// Slash format with time (e.g., 15/04/2025 13:45:00)
		"02/01/2006 15:04:05",
		// Slash date only (e.g., 15/04/2025)
		"02/01/2006",
		// Full timestamp without T (e.g., 2025-04-15 13:45:00)
		"2006-01-02 15:04:05",
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("cannot parse date '%s': %w", dateStr, err)
	}

	// Always return in the safe MSSQL datetime format (with time, zeroed if not provided)
	return t.Format("2006-01-02"), nil
}
// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`SELECT COUNT(1) FROM PriceDistUnitApproved WHERE priListCode = @p1 AND distUnitCode = @p2`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO PriceDistUnitApproved (
			priListCode, distUnitCode, distUnitName, fromDate, toDate
		) VALUES (
			@p1, @p2, @p3, @p4, @p5
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE PriceDistUnitApproved SET
			priListCode = @p1,
			distUnitCode = @p2,
			distUnitName = @p3,		
			fromDate = @p4,
			toDate = @p5
		WHERE   priListCode = @p1 AND distUnitCode = @p2
	`)
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec PriceDistUnitApproved) error {
	fixedFromDate, err := formatDateForMSSQL(rec.FromDate)
	fixedToDate, err1 := formatDateForMSSQL(rec.ToDate)
	if err != nil {
		return err
	}
	
	if err1 != nil {
		return err1
	}

	_, err = stmt.Exec(
		rec.PriListCode,				  
		rec.DistUnitCode,                   
		rec.DistUnitName,                 
		fixedFromDate,
		fixedToDate,
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec PriceDistUnitApproved) error {
	fixedFromDate, err := formatDateForMSSQL(rec.FromDate)
	fixedToDate, err1 := formatDateForMSSQL(rec.ToDate)
	if err != nil {
		return err
	}
	if err1 != nil {
		return err
	}

	_, err = stmt.Exec(
		rec.PriListCode,				  
		rec.DistUnitCode,                   
		rec.DistUnitName,                 
		fixedFromDate,
		fixedToDate,
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================

func SaveBatch(records []PriceDistUnitApproved,priListCode string) error {
	if len(records) == 0 {
		return nil
	}

	if db.DB == nil {
		return errors.New("database not initialized")
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statements
	checkStmt, err := prepareCheckStmt(tx)
	if err != nil {
		return err
	}
	defer checkStmt.Close()

	insertStmt, err := prepareInsertStmt(tx)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	updateStmt, err := prepareUpdateStmt(tx)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	// Process each record
	for _, rec := range records {

		rec.PriListCode = priListCode ;
		
		var count int
		if err = checkStmt.QueryRow(rec.PriListCode, rec.DistUnitCode).Scan(&count); err != nil {
			
			return err
		}
	
		if count > 0 {
			if err = execUpdate(updateStmt, rec); err != nil {
				
				return fmt.Errorf("update price entry failed for %s: %w", rec.DistUnitCode, err)
			}
		} else {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert price entry failed for %s: %w", rec.DistUnitCode, err)
			}
		}

		
		
	}

	return tx.Commit()
}