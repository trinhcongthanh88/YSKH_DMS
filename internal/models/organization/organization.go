package organization

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "YSKH_DMS/database"

	
	// "github.com/shopspring/decimal"
	// "github.com/davecgh/go-spew/spew"
)

// =====================
// ENTITY - MODEL Organization
// =====================

type Organization struct {
	// Các field từ API/JSON
	NodeCode              string          `json:"nodeCode"`
	NodeName              string          `json:"nodeName"`
	Status                string          `json:"status"`
	OrgType               string          `json:"orgType"`
	IsDistribute          bool            `json:"isDistribute"`
	CompCode              string          `json:"compCode"`
	CompName              string          `json:"compName"`
	CreateUser            string          `json:"createUser"`
	CreateDate            string          `json:"createDate"`

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

	var t time.Time
	var err error

	if strings.Contains(dateStr, "T") {
		dateStr = strings.Replace(dateStr, "T", " ", 1)
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
	} else if strings.Contains(dateStr, "/") {
		t, err = time.Parse("02/01/2006 15:04:05", dateStr)
	} else {
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
	}

	if err != nil {
		return "", fmt.Errorf("cannot parse date '%s': %w", dateStr, err)
	}

	return t.Format("2006-01-02 15:04:05"), nil
}

// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`SELECT COUNT(1) FROM Organization WHERE nodeCode = @p1`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO Organization (
			nodeCode, nodeName, status, orgType, isDistribute, compCode, compName, createUser, createDate
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE Organization SET
			nodeCode = @p1,                -- Cập nhật nodeCode bằng nodeCode
			nodeName = @p2,
			status = @p3,
			orgType = @p4,
			isDistribute = @p5,
			compCode = @p6,
			compName = @p7,
			createUser = @p8,
			createDate = @p9
		WHERE nodeCode = @p1
	`)	
		
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec Organization) error {
	fixedDate, err := formatDateForMSSQL(rec.CreateDate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		rec.NodeCode,                   // @p1
		rec.NodeName,                   // @p2
		rec.Status,                     // @p3
		rec.OrgType,                    // @p4
		rec.IsDistribute,               // @p5
		rec.CompCode,                   // @p6	
		rec.CompName,                   // @p7
		rec.CreateUser,                 // @p8
		fixedDate,                      // @p9
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec Organization) error {
	fixedDate, err := formatDateForMSSQL(rec.CreateDate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		rec.NodeCode,                   // @p1
		rec.NodeName,                   // @p2
		rec.Status,                     // @p3
		rec.OrgType,                    // @p4
		rec.IsDistribute,               // @p5
		rec.CompCode,                   // @p6	
		rec.CompName,                   // @p7
		rec.CreateUser,                 // @p8
		fixedDate,                      // @p9
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================
func GetAllOrganization() ([]Organization, error) {
	if db.DB == nil {
		return nil, errors.New("database not initialized")
	}	
	rows, err := db.DB.Query(`SELECT
		nodeCode, nodeName
		FROM Organization`)
	if  err != nil {
		return nil, err
	}	
	defer rows.Close()
	var rowsData []Organization
	for rows.Next() {
		var rec Organization								
		if err := rows.Scan(
			&rec.NodeCode,
			&rec.NodeName,
		); err != nil {
			return nil, err
		}
		rowsData =  append(rowsData, rec)
	}
	
	return rowsData, nil
}
func SaveBatchApiWeb(records []Organization) error {
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
	// Process each record
	for _, rec := range records {
		
		var count int
		if err = checkStmt.QueryRow(rec.NodeCode).Scan(&count); err != nil {
			return err
		}

		if count < 0 {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert failed for %s: %w", rec.NodeCode, err)
			}
		}
	}
	return tx.Commit()
}
func SaveBatch(records []Organization) error {
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
	updateStmt, err := prepareUpdateStmt(tx)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	// Process each record
	for _, rec := range records {
	
		if err = execUpdate(updateStmt, rec); err != nil {
			return fmt.Errorf("update Product UOM failed for %s: %w", rec.NodeCode, err)
		}
	}

	return tx.Commit()
}