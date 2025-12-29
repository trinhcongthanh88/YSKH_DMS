package stafflist

import (
	db "YSKH_DMS/database"
	"database/sql"
	"errors"

	// "fmt"
	// "os"
	base "YSKH_DMS/internal/models"

	"github.com/davecgh/go-spew/spew"
)

type Staff struct {
	StaffCode     string     `json:"staffCode"`
	StaffName     string     `json:"staffName"`
	Account       string     `json:"account"`
	Status        string     `json:"status"`
	StaffJobTitle []staffJob `json:"staffJobTitle"`
	CreateDate    string     `json:"createDate"`
	CreateUser    string     `json:"createUser"`
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

func execUpdate(stmt *sql.Stmt, rec Staff) error {

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

func execInsert(stmt *sql.Stmt, rec Staff) error {
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

func SaveBatch(records []Staff) error {

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
