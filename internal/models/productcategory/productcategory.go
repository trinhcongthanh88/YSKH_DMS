package productcategory

import (
	"database/sql"
	"errors"
    "fmt"
    "strings"
    "time"
	db  "YSKH_DMS/database"
	// "fmt"
	// "os"
	"github.com/davecgh/go-spew/spew"
)


// =====================
// ENTITY
// =====================

type ProductCategory struct {
	ProCategoryCode       string `json:"proCategoryCode"`
    ProCategoryName       string `json:"proCategoryName"`
    ProCategoryParentCode string `json:"proCategoryParentCode"`
    ProCategoryCreateDate string `json:"proCategoryCreateDate"`
    ProCategoryCreateName string `json:"proCategoryCreateName"`
 
}
func formatDateForMSSQL(dateStr string) (any, error) {
    dateStr = strings.TrimSpace(dateStr)
    if dateStr == "" || 
	   strings.EqualFold(dateStr, "null") || 
	   strings.EqualFold(dateStr, "nil") {
		return nil, nil // No error, just no value
	}
    var t time.Time
    var err error

    // Thử các định dạng phổ biến từ API
    if strings.Contains(dateStr, "T") {
        // ISO: 2025-01-08T16:30:23
        dateStr = strings.Replace(dateStr, "T", " ", 1)
        t, err = time.Parse("2006-01-02 15:04:05", dateStr)
    } else if strings.Contains(dateStr, "/") {
        // DD/MM/YYYY HH:MM:SS
        t, err = time.Parse("02/01/2006 15:04:05", dateStr)
    } else {
        // Fallback
        t, err = time.Parse("2006-01-02 15:04:05", dateStr)
    }

    if err != nil {
        return "", fmt.Errorf("cannot parse date '%s': %w", dateStr, err)
    }

    // Trả về định dạng an toàn nhất cho MSSQL
    return t.Format("2006-01-02 15:04:05"), nil
}

// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		SELECT COUNT(1)
        FROM ProductCategory
        WHERE proCategoryCode = @p1 
	`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	
	return tx.Prepare(`
		INSERT INTO ProductCategory (
			proCategoryCode,
			proCategoryName,
			proCategoryParentCode,
			proCategoryCreateDate,
			proCategoryCreateName
		
		)
		VALUES (@p1,@p2,@p3,@p4,@p5)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE ProductCategory
		SET
			proCategoryName = @p2?,
			proCategoryParentCode = @p3,
			proCategoryCreateDate = @p4,
			proCategoryCreateName =@p5
		WHERE proCategoryCode = @p1 
	`)
}


// =====================
// EXEC HELPERS
// =====================

func execUpdate(stmt *sql.Stmt, rec ProductCategory) error {
	fixedDate, err := formatDateForMSSQL(rec.ProCategoryCreateDate)
		if err != nil {
			return err
		}
	_, err1 := stmt.Exec(
		rec.ProCategoryCode,
		rec.ProCategoryName,
		rec.ProCategoryParentCode,
		fixedDate,
		rec.ProCategoryCreateName,
	)
	return err1
}

func execInsert(stmt *sql.Stmt, rec ProductCategory) error {
	fixedDate, err := formatDateForMSSQL(rec.ProCategoryCreateDate)
		if err != nil {
			return err
		}
	_, err1 := stmt.Exec(
		rec.ProCategoryCode,
		rec.ProCategoryName,
		rec.ProCategoryParentCode,
		fixedDate,
		rec.ProCategoryCreateName,
	)
		

	return err1
}


// =====================
// SAVE BATCH (UPSERT)
// =====================

func SaveBatch(records []ProductCategory) error {

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

	
	for _, rec := range records {
		var count int
		err := checkStmt.
			QueryRow(rec.ProCategoryCode).
			Scan(&count)
		

		if err != nil {
			return err
		}
		
		if count > 0 {
			err = execUpdate(updateStmt, rec)
		} else {
			err = execInsert(insertStmt, rec)
				 spew.Dump("sssssssss")
		spew.Dump(err)
		}
	   
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}