package product

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "YSKH_DMS/database"
	producttaxModel "YSKH_DMS/internal/models/producttax"
	productuomModel "YSKH_DMS/internal/models/productuom"
	
	"github.com/shopspring/decimal"
	"github.com/davecgh/go-spew/spew"
)

// =====================
// ENTITY - MODEL PRODUCT
// =====================

type Product struct {
	// productId không lấy từ JSON, tự động gán bằng proCode
	ProductId string `json:"-"`

	// Các field từ API/JSON
	ProCode                string          `json:"proCode"`
	ProName                string          `json:"proName"`
	ProType                string          `json:"proType"`
	ProStatus              bool            `json:"proStatus"`
	ProIsSale              bool            `json:"proIsSale"`
	ProIsVisible           bool            `json:"proIsVisible"`
	ProIsVar               bool            `json:"proIsVar"`
	ProManagerLot          bool            `json:"proManagerLot"`
	ProManagerBarcode      bool            `json:"proManagerBarcode"`
	ProTaxGroup            string          `json:"proTaxGroup"`
	TaxRate                decimal.Decimal `json:"taxRate"`
	PriceSellDefaultNotVat decimal.Decimal `json:"priceSellDefaultNotVat"`
	PriceSelDefaultVat     decimal.Decimal `json:"priceSelDefaultVat"`
	ProCategoryCode        string          `json:"proCategoryCode"`
	ProCategoryName        string          `json:"proCategoryName"`
	ProCategoryParentCode  string          `json:"proCategoryParentCode"`
	ProCreateDate          string          `json:"proCreateDate"`
	ProCreateName          string          `json:"proCreateName"`
	ProUomId	 		   string          `json:"proUomId"`
	ProTaxId	 		   string          `json:"proTaxId"`
	ProUomDefault *productuomModel.ProductUOM `json:"proUomDefault,omitempty"` // Có thể null
	ProCategory []ProCategoryItem `json:"proCategory"`
	ProListTax []producttaxModel.ProductTax `json:"proListTax"`
	ProListUom []productuomModel.ProductUOM `json:"proListUom"`
}

type ProCategoryItem struct {
	ProCategoryCode        string  `json:"proCategoryCode"`        
	ProCategoryName        string  `json:"proCategoryName"`       
	ProCategoryParentCode  string  `json:"proCategoryParentCode"`   
	
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
	return tx.Prepare(`SELECT COUNT(1) FROM Product WHERE proCode = @p1`)
}

func prepareInsertStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		INSERT INTO Product (
			productId, proCode, proName, proType, proStatus, proIsSale, proIsVisible,
			proIsVar, proManagerLot, proManagerBarcode, proTaxGroup, taxRate,
			priceSellDefaultNotVat, priceSelDefaultVat, proCategoryCode,
			proCategoryName, proCategoryParentCode, proCreateDate, proCreateName,proUomId,proTaxId
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10,
			@p11, @p12, @p13, @p14, @p15, @p16, @p17, @p18, @p19,@p20,@p21
		)
	`)
}

func prepareUpdateStmt(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(`
		UPDATE Product SET
			productId = @p1,                -- Cập nhật productId bằng proCode
			proName = @p2,
			proType = @p3,
			proStatus = @p4,
			proIsSale = @p5,
			proIsVisible = @p6,
			proIsVar = @p7,
			proManagerLot = @p8,
			proManagerBarcode = @p9,
			proTaxGroup = @p10,
			taxRate = @p11,
			priceSellDefaultNotVat = @p12,
			priceSelDefaultVat = @p13,
			proCategoryCode = @p14,
			proCategoryName = @p15,
			proCategoryParentCode = @p16,
			proCreateDate = @p17,
			proCreateName = @p18,
			proUomId = @p20,
			proTaxId = @p21
		WHERE proCode = @p19
	`)
}

// =====================
// EXEC: INSERT & UPDATE
// =====================

func execInsert(stmt *sql.Stmt, rec Product) error {
	fixedDate, err := formatDateForMSSQL(rec.ProCreateDate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		rec.ProductId,                    // @p1  → ProductId 
		rec.ProCode,                    // @p2  → proCode
		rec.ProName,                    // @p3
		rec.ProType,                    // @p4
		rec.ProStatus,                  // @p5
		rec.ProIsSale,                  // @p6
		rec.ProIsVisible,               // @p7
		rec.ProIsVar,                   // @p8
		rec.ProManagerLot,              // @p9
		rec.ProManagerBarcode,          // @p10
		rec.ProTaxGroup,                // @p11
		rec.TaxRate,                    // @p12
		rec.PriceSellDefaultNotVat,     // @p13
		rec.PriceSelDefaultVat,         // @p14
		rec.ProCategoryCode,            // @p15
		rec.ProCategoryName,            // @p16
		rec.ProCategoryParentCode,      // @p17
		fixedDate,                      // @p18
		rec.ProCreateName,              // @p19
		rec.ProUomId, 					// @p20
		rec.ProTaxId, 					// @p21
	
	)
	return err
}

func execUpdate(stmt *sql.Stmt, rec Product) error {
	fixedDate, err := formatDateForMSSQL(rec.ProCreateDate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		rec.ProductId,                   // @p1  → productId =
		rec.ProName,                    // @p2
		rec.ProType,                    // @p3
		rec.ProStatus,                  // @p4
		rec.ProIsSale,                  // @p5
		rec.ProIsVisible,               // @p6
		rec.ProIsVar,                   // @p7
		rec.ProManagerLot,              // @p8
		rec.ProManagerBarcode,          // @p9
		rec.ProTaxGroup,                // @p10
		rec.TaxRate,                    // @p11
		rec.PriceSellDefaultNotVat,     // @p12
		rec.PriceSelDefaultVat,         // @p13
		rec.ProCategoryCode,            // @p14
		rec.ProCategoryName,            // @p15
		rec.ProCategoryParentCode,      // @p16
		fixedDate,                      // @p17
		rec.ProCreateName,              // @p18
		rec.ProCode,                    // @p19 → WHERE
		rec.ProUomId, 					// @p20
		rec.ProTaxId, 					// @p21
	)
	return err
}

// =====================
// SAVE BATCH - UPSERT
// =====================

func SaveBatch(records []Product) error {
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
		
		if rec.ProUomDefault != nil {
			defaultUom := *rec.ProUomDefault // copy struct

			// Nếu ProListUom rỗng hoặc nil → khởi tạo với default
			if len(rec.ProListUom) == 0 {
				rec.ProListUom = []productuomModel.ProductUOM{defaultUom}
			} else {
				// Nếu đã có → thêm default vào đầu danh sách
				rec.ProListUom = append([]productuomModel.ProductUOM{defaultUom}, rec.ProListUom...)
			}
			if rec.ProUomDefault.ProductId != "" {
				rec.ProductId = rec.ProUomDefault.ProductId
			}
		}
		
		if len(rec.ProListTax) > 0 && rec.ProListTax[0].Id != "" {
				rec.ProTaxId = rec.ProListTax[0].Id
		}
		if rec.ProCategory != nil && len(rec.ProCategory) > 0 {
			firstCat := rec.ProCategory[0]
			rec.ProCategoryCode = firstCat.ProCategoryCode
			rec.ProCategoryName = firstCat.ProCategoryName
			rec.ProCategoryParentCode = firstCat.ProCategoryParentCode
		} 
		
		spew.Dump(rec);
		var count int
		if err = checkStmt.QueryRow(rec.ProCode).Scan(&count); err != nil {
			return err
		}

		if count > 0 {
			if err = execUpdate(updateStmt, rec); err != nil {
				return fmt.Errorf("update failed for %s: %w", rec.ProCode, err)
			}
			
		} else {
			if err = execInsert(insertStmt, rec); err != nil {
				return fmt.Errorf("insert failed for %s: %w", rec.ProCode, err)
			}
		}
		if len(rec.ProListTax) > 0{
			producttaxModel.SaveBatch(rec.ProListTax) ;
		}
		
		if len(rec.ProListUom) > 0{
			productuomModel.SaveBatch(rec.ProListUom) ;
		}
		if len(rec.ProListTax) > 0{
			producttaxModel.SaveBatch(rec.ProListTax) ;
		}

		
		
	}

	return tx.Commit()
}