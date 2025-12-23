package productuom

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

type ProductUOM struct {
	Id 					  string 		   `json:"id"`
	ProductId             string          `json:"productId"`
	Code                  string          `json:"code"`
	Title                 string          `json:"title"`
	IsDefault             bool            `json:"isDefault"`
	Ratio                 decimal.Decimal `json:"ratio"`
	VisualIndex           int             `json:"visualIndex,omitempty"`
	
}

// =====================
// HELPER: FORMAT DATE CHO MSSQL
// =====================



// =====================
// PREPARE STATEMENTS
// =====================

func prepareCheckStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`SELECT COUNT(1) FROM ProductUOM WHERE id = @p1`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO ProductUOM (
			id, productId, code, title, isDefault, ratio, visualIndex
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE ProductUOM SET
			productId = @p2,
			code = @p3,
			title = @p4,
			isDefault = @p5,
			ratio = @p6,
			visualIndex = @p7,
		WHERE id = @p1
	`)
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec ProductUOM) error {
	_, err := stmt.Exec(
		rec.Id,                     
		rec.ProductId,                  
		rec.Code,                   
		rec.Title,                 
		rec.IsDefault,             
		rec.Ratio,  
		rec.VisualIndex,  
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec ProductUOM) error {

	_, err := stmt.Exec(
	    rec.Id,                     
		rec.ProductId,                  
		rec.Code,                   
		rec.Title,                 
		rec.IsDefault,             
		rec.Ratio,  
		rec.VisualIndex,  
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================

func SaveBatch(records []ProductUOM) error {
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
				return fmt.Errorf("update Product UOM failed for %s: %w", rec.Id, err)
			}
		} else {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert Product UOM failed for %s: %w", rec.Id, err)
			}
		}

		
		
	}

	return tx.Commit()
}