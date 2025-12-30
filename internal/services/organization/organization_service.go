package product

import (
	viettelService "YSKH_DMS/internal/services/viettel"
	// "fmt"
	"encoding/json"
	// "os"
	"log"
	"strings"
	"time"
	organizationModel "YSKH_DMS/internal/models/organization"
	// "github.com/davecgh/go-spew/spew"
	
)

type ViettelResponse struct {
	Status StatusInfo                 `json:"status"`
	Data   []organizationModel.Organization     `json:"data"`
	Pagination PagingInfo             `json:"pagination"`
}
type ViettelResponseWeb struct {
	Content   []organizationModel.Organization     `json:"content"`
	TotalPages int             `json:"totalPages"`
	TotalElements int   `json:"totalElements"`
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

func SaveItemDmsViettel() (any, error) {
	// Lấy danh sách tổ chức
	itemOrganizations, err := organizationModel.GetAllOrganization()
	if err != nil {
		return nil, err
	}

	if len(itemOrganizations) == 0 {
		log.Println("Không có tổ chức nào để đồng bộ từ Viettel DMS")
		return 0, nil
	}

	totalSaved := 0

	for _, org := range itemOrganizations {
		log.Printf("Đang xử lý tổ chức: %s - %s (NodeCode: %s)\n",
			org.NodeName, org.NodeCode, org.NodeCode)

		// Query params
		queryApi := map[string]any{
			"distCode":    org.NodeCode,
			"distName":    org.NodeName,
			"status":      []string{}, // Nếu API chấp nhận mảng rỗng = lấy tất cả trạng thái
			"createDate":  "2025-01-21",
		}

		urlapi := "https://app.vietteldms.com/openapi/v1/GetListDist"
		page := 0
		size := 100

		// Thử gọi API với retry khi bị rate limit
		var repData []byte
		var result ViettelResponse
		success := false

		for retry := 0; retry < 3; retry++ { // Thử tối đa 3 lần
			urlApiBuild, err := viettelService.BuildViettelDistURL(urlapi, queryApi, page, size)
			if err != nil {
				log.Printf("Lỗi build URL cho org %s: %v", org.NodeCode, err)
				break
			}

			repData, err = viettelService.ViettelGet(urlApiBuild)
			if err != nil {
				log.Printf("Lỗi HTTP cho org %s: %v", org.NodeCode, err)
				break
			}

			// Debug response (có thể comment lại khi chạy production)
			// spew.Dump(string(repData))

			// Kiểm tra rate limit
			if strings.Contains(string(repData), "Rate limit exceeded") {
				waitTime := (retry + 1) * 60 // 60s, 120s, 180s
				log.Printf("Rate limit exceeded cho org %s. Chờ %d giây rồi thử lại... (lần %d/3)\n",
					org.NodeCode, waitTime, retry+1)
				time.Sleep(time.Duration(waitTime) * time.Second)
				continue
			}

			// Parse JSON
			if err := json.Unmarshal(repData, &result); err != nil {
				log.Printf("Lỗi parse JSON cho org %s: %v\nRaw: %s", org.NodeCode, err, string(repData))
				break
			}

			// spew.Dump(result) // Debug struct nếu cần

			success = true
			break
		}

		if !success {
			log.Printf("Bỏ qua org %s sau 3 lần thử\n", org.NodeCode)
			// Delay thêm để tránh làm nặng server
			time.Sleep(5 * time.Second)
			continue
		}

		// Lưu dữ liệu
		if len(result.Data) > 0 {
			if err := organizationModel.SaveBatch(result.Data); err != nil {
				log.Printf("Lỗi lưu batch cho org %s: %v", org.NodeCode, err)
			} else {
				log.Printf("Đã lưu %d bản ghi cho org %s\n", len(result.Data), org.NodeCode)
				totalSaved += len(result.Data)
			}
		} else {
			log.Printf("Không có dữ liệu mới cho org %s\n", org.NodeCode)
		}

		// Delay giữa các org để tránh rate limit (rất quan trọng!)
		time.Sleep(5 * time.Second) // Điều chỉnh: 5-15s tùy số lượng org
	}

	log.Printf("Hoàn thành đồng bộ Viettel DMS. Tổng cộng lưu: %d bản ghi\n", totalSaved)
	return totalSaved, nil
}
func SaveBatchDmsViettelWeb() (any, error) {


	urlapi := "https://app.vietteldms.vn/api/v1/dms-service/distUnit"
	page := 0 
	size := 20 

	for {
		urlApiBuild, err := viettelService.BuildViettelDistURLWeb(urlapi,page,size)
		if err != nil {
			return nil, err
		}

		repData, err := viettelService.ViettelGetApiweb(urlApiBuild)
		if err != nil {
			return nil, err
		}
	
		var result ViettelResponseWeb
		if err := json.Unmarshal(repData, &result); err != nil {
			
			return nil, err
		}
	
		
		err = organizationModel.SaveBatchApiWeb(result.Content)
		if page >= result.TotalPages-1 {
			break // Dừng khi là trang cuối
		}
	
		
		

		page++ // Tăng trang
	}
	
    
    return page, nil

}