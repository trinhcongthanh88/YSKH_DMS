package priceentry

import (
	"database/sql"
	"errors"
	"fmt"
	_"strings"
	_"time"

	db "YSKH_DMS/database"
	"github.com/shopspring/decimal"
	// "github.com/davecgh/go-spew/spew"
)

// =====================
// ENTITY - MODEL PriceEntry
// =====================

type PriceEntry struct {
	PriceEntryId 		   int 		  	   `json:"-"`
	PriListCode            string      	   `json:"priListCode"`
	ProdCode               string          `json:"prodCode"`
    ProdName               string          `json:"prodName"`
	ProdTaxClass           string          `json:"prodTaxClass"`
	ProPriceNotVat		   decimal.Decimal `json:"proPriceNotVat"`
	ProPriceVat		       decimal.Decimal `json:"proPriceVat"`
	
}
// =====================
// HELPER: FORMAT DATE CHO MSSQL
// =====================



// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`SELECT COUNT(1) FROM PriceEntry WHERE priListCode = @p1 AND prodCode = @p2`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO PriceEntry (
			 priListCode, prodCode, prodName, prodTaxClass, proPriceNotVat, proPriceVat
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE PriceEntry SET
			priListCode = @p1,
			prodCode = @p2,
			prodName = @p3,
			prodTaxClass = @p4,
			proPriceNotVat = @p5,
			proPriceVat = @p6
		WHERE  prodCode = @p2 AND priListCode = @p1	
	`)
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec PriceEntry) error {
	_, err := stmt.Exec(
		rec.PriListCode,                  
		rec.ProdCode,                   
		rec.ProdName,                 
		rec.ProdTaxClass,             
		rec.ProPriceNotVat,
		rec.ProPriceVat,
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec PriceEntry) error {

	_, err := stmt.Exec(
	   rec.PriListCode,                  
		rec.ProdCode,                   
		rec.ProdName,                 
		rec.ProdTaxClass,             
		rec.ProPriceNotVat,
		rec.ProPriceVat,
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================

func SaveBatch(records []PriceEntry, priListCode string) error {
	  
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
		if err = checkStmt.QueryRow(rec.PriListCode, rec.ProdCode).Scan(&count); err != nil {
				
			return err
		}

		if count > 0 {
			if err = execUpdate(updateStmt, rec); err != nil {
				return fmt.Errorf("update price entry failed for %s: %w", rec.PriceEntryId, err)
			}
		} else {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert price entry failed for %s: %w", rec.PriceEntryId, err)
			}
		}

		
		
	}

	return tx.Commit()
}