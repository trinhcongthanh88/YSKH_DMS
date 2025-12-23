package producttax

import (
	"database/sql"
	"errors"
	"fmt"
	_"strings"
	_"time"

	db "YSKH_DMS/database"
	"github.com/shopspring/decimal"
	"github.com/davecgh/go-spew/spew"
)

// =====================
// ENTITY - MODEL PRODUCT
// =====================

type ProductTax struct {
	Id 					   string 		   `json:"id"`
	TaxClassId             string          `json:"taxClassId"`
	TaxRate                decimal.Decimal `json:"taxRate"`
    Piggyback              bool            `json:"piggyback"`
	TaxPriority            int             `json:"taxPriority,omitempty"`
	
}
// =====================
// HELPER: FORMAT DATE CHO MSSQL
// =====================



// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`SELECT COUNT(1) FROM ProductTax WHERE id = @p1`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO ProductTax (
			id, taxClassId, taxRate, piggyback, taxPriority
		) VALUES (
			@p1, @p2, @p3, @p4, @p5
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE ProductTax SET
			taxClassId = @p2,
			taxRate = @p3,
			piggyback = @p4,
			taxPriority = @p5,
		WHERE id = @p1
	`)
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec ProductTax) error {
	_, err := stmt.Exec(
		rec.Id,                     
		rec.TaxClassId,                  
		rec.TaxRate,                   
		rec.Piggyback,                 
		rec.TaxPriority,             
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec ProductTax) error {

	_, err := stmt.Exec(
	    rec.Id,                     
		rec.TaxClassId,                  
		rec.TaxRate,                   
		rec.Piggyback,                 
		rec.TaxPriority,    
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================

func SaveBatch(records []ProductTax) error {
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
		
	
		spew.Dump(rec);
		var count int
		if err = checkStmt.QueryRow(rec.Id).Scan(&count); err != nil {
			return err
		}

		if count > 0 {
			if err = execUpdate(updateStmt, rec); err != nil {
				return fmt.Errorf("update Product tax failed for %s: %w", rec.Id, err)
			}
		} else {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert Product tax failed for %s: %w", rec.Id, err)
			}
		}

		
		
	}

	return tx.Commit()
}