package stafflist

import (
	viettelService "YSKH_DMS/internal/services/viettel"
	// "fmt"
	"encoding/json"
	// "os"
	response "YSKH_DMS/internal/models"
	staffResponse "YSKH_DMS/internal/models/stafflist"
)

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

		var result response.ResponseGenericList[staffResponse.Staff]
		if err := json.Unmarshal(repData, &result); err != nil {

			return nil, err
		}

		err = staffResponse.SaveBatch(result.List)
		if page >= result.Pagination.TotalPages-1 {
			break // Dừng khi là trang cuối
		}

		page++ // Tăng trang
	}

	return page, nil

}
