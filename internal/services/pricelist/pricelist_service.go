package pricelist

import (
	viettelService "YSKH_DMS/internal/services/viettel"
	// "fmt"
	"encoding/json"
	// "os"
	priceListModel "YSKH_DMS/internal/models/pricelist"
	// "github.com/davecgh/go-spew/spew"
	
)

type ViettelResponse struct {
	Status StatusInfo                 `json:"status"`
	Data   []priceListModel.PriceList     `json:"data"`
	Pagination PagingInfo             `json:"pagination"`
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
		// "priListCode":"",
		// "priListName":"",
		"priListStatus":"ACT",
		// "priListType":"SELL"
	}
	

	urlapi := "https://app.vietteldms.com/openapi/v1/GetListPrice"
	page := 0 
	size := 100 

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
	
		
		err = priceListModel.SaveBatch(result.Data)
		if page >= result.Pagination.TotalPages-1 {
			break // Dừng khi là trang cuối
		}
	
		
		

		page++ // Tăng trang
	}
	
    
    return page, nil

}