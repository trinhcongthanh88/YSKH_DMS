package stafflist

import (
	db "YSKH_DMS/database"
	base "YSKH_DMS/internal/models"
	"database/sql"
	"errors"

	// "fmt"
	// "os"
	"github.com/davecgh/go-spew/spew"
)

// "fmt"
// "os"

type staffJob struct {
	FromDate     string `json:"fromDate"`
	DistUnitCode string `json:"distUnitCode"`
	DistUnitName string `json:"distUnitName"`
	JobTitleCode string `json:"jobTitleCode"`
	JobTitleName string `json:"jobTitleName"`
}

// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmtStaffJob(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		SELECT COUNT(1)
        FROM ProductCategory
        WHERE proCategoryCode = @p1 
	`)
}

func prepareInsertStmtStaffJob(tx *sql.Tx) (*sql.Stmt, error) {

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

func prepareUpdateStmtStaffJob(tx *sql.Tx) (*sql.Stmt, error) {
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

func execUpdateStaffJob(stmt *sql.Stmt, rec Staff) error {
	fixedDate, err := base.FormatDateForMSSQL(rec.CreateDate)
	if err != nil {
		return err
	}
	_, err1 := stmt.Exec(
		rec.StaffCode,
		rec.StaffName,
		rec.Account,
		fixedDate,
		rec.Status,
		// rec.StaffJobTitle
	)
	return err1
}

func execInsertStaffJob(stmt *sql.Stmt, rec Staff) error {
	fixedDate, err := base.FormatDateForMSSQL(rec.CreateDate)
	if err != nil {
		return err
	}
	_, err1 := stmt.Exec(
		rec.StaffCode,
		rec.StaffName,
		rec.Account,
		fixedDate,
		rec.Status,
		// rec.StaffJobTitle
	)

	return err1
}

// =====================
// SAVE BATCH (UPSERT)
// =====================

func SaveBatchStaffJob(records []Staff) error {

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
			QueryRow(rec.StaffCode).
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
