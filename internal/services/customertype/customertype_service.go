package customertype

import (
	viettelService "YSKH_DMS/internal/services/viettel"
	// "fmt"
	"encoding/json"
	// "os"
	cutomertypeModel "YSKH_DMS/internal/models/customertype"
	// "github.com/davecgh/go-spew/spew"
	
)

type ViettelResponse struct {
	Status StatusInfo                           `json:"status"`
	Data   []cutomertypeModel.CustomerType     `json:"data"`
	Pagination PagingInfo                           `json:"pagination"`
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

	queryApi := map[string]any{
		"status":  "ACT",
	}

	urlapi := "https://app.vietteldms.com/openapi/v1/GetCustomerTypeList"
	page := 0 
	size := 20 

	for {
		urlApiBuild, err := viettelService.BuildViettelDistURL(urlapi, queryApi,page,size)
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
		
		if page >= result.Pagination.TotalPages-1 {
			break // Dừng khi là trang cuối
		}
		
		err = cutomertypeModel.SaveBatch(result.Data)
		

		page++ // Tăng trang
	}
	
    
    return page, nil

}