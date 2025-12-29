package productcategory

import (
	viettelService "YSKH_DMS/internal/services/viettel"
	// "fmt"
	"encoding/json"
	// "os"
	productcategoryModel "YSKH_DMS/internal/models/productcategory"
)

type ViettelResponse struct {
	Status     StatusInfo                             `json:"status"`
	Data       []productcategoryModel.ProductCategory `json:"data"`
	Pagination PagingInfo                             `json:"pagination"`
}

type StatusInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PagingInfo struct {
	Last             bool `json:"last"`
	First            bool `json:"first"`
	Empty            bool `json:"empty"`
	NumberOfElements int  `json:"numberOfElements"`
	TotalElements    int  `json:"totalElements"`
	TotalPages       int  `json:"totalPages"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
}

func SaveBatchDmsViettel() (any, error) {

	queryApi := map[string]any{}

	urlapi := "https://app.vietteldms.com/openapi/v1/GetAllProductCategory"
	page := 0
	size := 20

	for {
		urlApiBuild, err := viettelService.BuildViettelDistURL(urlapi, queryApi, page, size)
		if err != nil {
			return nil, err
		}

		repData, err := viettelService.ViettelGet(urlApiBuild)
		if err != nil {
			return nil, err
		}

		var result ViettelResponse
		if err := json.Unmarshal(repData, &result); err != nil {

			return nil, err
		}

		err = productcategoryModel.SaveBatch(result.Data)
		if page >= result.Pagination.TotalPages-1 {
			break // Dừng khi là trang cuối
		}

		page++ // Tăng trang
	}

	return page, nil

}
