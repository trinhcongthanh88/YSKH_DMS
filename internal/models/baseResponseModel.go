package models

import (
	"fmt"
	"strings"
	"time"
)

// Pagination đại diện cho cấu trúc "pagination"
type Pagination struct {
	Last             bool `json:"last"`
	First            bool `json:"first"`
	Empty            bool `json:"empty"`
	NumberOfElements int  `json:"numberOfElements"`
	TotalElements    int  `json:"totalElements"`
	TotalPages       int  `json:"totalPages"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

// ResponseGeneric là cấu trúc model generic
// T là type parameter, sẽ được thay thế bằng kiểu dữ liệu thực tế của trường 'data'
type ResponseGeneric[T any] struct {
	Status     StatusInfo `json:"status"`
	Data       T          `json:"data"` // Trường Data sử dụng kiểu generic T
	Pagination Pagination `json:"pagination"`
}
type ResponseGenericList[T any] struct {
	Status     StatusInfo `json:"status"`
	List       []T        `json:"data"` // Trường Data sử dụng kiểu generic T
	Pagination Pagination `json:"pagination"`
}

type StatusInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func FormatDateForMSSQL(dateStr string) (any, error) {
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
