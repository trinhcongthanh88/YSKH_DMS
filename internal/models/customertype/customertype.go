package customertype

import (
	"database/sql"
	"errors"

	db  "YSKH_DMS/database"
	// "fmt"
	// "os"
	// "github.com/davecgh/go-spew/spew"
)


// =====================
// ENTITY
// =====================

type CustomerType struct {
	GroupCode      string `json:"groupCode"`
    GroupName      string `json:"groupName"`
    GroupStatus    string `json:"groupStatus"`
    CustTypeCode   string `json:"custTypeCode"`
    CustTypeName   string `json:"custTypeName"`
    CustTypeStatus string `json:"custTypeStatus"`
    CustTypeUser   string `json:"custTypeUser"`
    CustTypeDate   string `json:"custTypeDate"`
}


// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		SELECT COUNT(1)
        FROM CustomerType
        WHERE custTypeCode = @p1 
        AND groupCode = @p2
	`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO CustomerType (
			custTypeCode,
			groupCode,
			groupName,
			groupStatus,
			custTypeName,
			custTypeStatus,
			custTypeUser,
			custTypeDate
		
		)
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE CustomerType
		SET
			groupName = @p1?,
			groupStatus = @p2,
			custTypeName = @p3,
			custTypeStatus =@p4,
			custTypeUser = @p5,
			custTypeDate = @p6,
		WHERE custTypeCode = ? AND groupCode = ?
	`)
}


// =====================
// EXEC HELPERS
// =====================

func execUpdate(stmt *sql.Stmt, rec CustomerType) error {
	_, err := stmt.Exec(
		rec.GroupName,
		rec.GroupStatus,
		rec.CustTypeName,
		rec.CustTypeStatus,
		rec.CustTypeUser,
		rec.CustTypeDate,
		rec.CustTypeCode,
		rec.GroupCode,
	)
	return err
}

func execInsert(stmt *sql.Stmt, rec CustomerType) error {
	_, err := stmt.Exec(
		rec.CustTypeCode,
		rec.GroupCode,
		rec.GroupName,
		rec.GroupStatus,
		rec.CustTypeName,
		rec.CustTypeStatus,
		rec.CustTypeUser,
		rec.CustTypeDate,
		
	)
		

	return err
}


// =====================
// SAVE BATCH (UPSERT)
// =====================

func SaveBatch(records []CustomerType) error {

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
			QueryRow(rec.CustTypeCode, rec.GroupCode).
			Scan(&count)
		
	
		if err != nil {
			return err
		}
		
		if count > 0 {
			err = execUpdate(updateStmt, rec)
		} else {
			err = execInsert(insertStmt, rec)
		
		}

		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}