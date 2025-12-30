package dms

import (
	_ "bytes"
	_ "encoding/json"
	_ "io"
	_ "net/http"
	_ "net/url"    
	_ "strings"   
	_ "time"
	_ "encoding/json"
	customertypeService "YSKH_DMS/internal/services/customertype"
	productcategoryService "YSKH_DMS/internal/services/productcategory"
	productService "YSKH_DMS/internal/services/product"
	pricelistService "YSKH_DMS/internal/services/pricelist"
	organizationService "YSKH_DMS/internal/services/organization"
	"github.com/gin-gonic/gin"
	
)

// // GET /users: Lấy danh sách users từ DB
// func GetUsers(c *gin.Context) {
// 	users, err := models.GetUsers()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
func CustomertypeApi(c *gin.Context) {

	numberPage, err := customertypeService.SaveBatchDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}
func ProductcategoryApi(c *gin.Context) {

	numberPage, err := productcategoryService.SaveBatchDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}
func ProductApi(c *gin.Context) {

	numberPage, err := productService.SaveBatchDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}
func PriceListApi(c *gin.Context) {

	numberPage, err := pricelistService.SaveBatchDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}
func OrganizationListApiWeb(c *gin.Context) {

	numberPage, err := organizationService.SaveBatchDmsViettelWeb();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}
func OrganizationListApi(c *gin.Context) {

	numberPage, err := organizationService.SaveItemDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"page":    numberPage,
	})

}


// // POST /users: Tạo user mới và lưu vào DB
// func CreateUser(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if err := models.CreateUser(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, user)
// }

// // GET /external-api: Kết nối đến một API external (ví dụ GET)
// func GetExternalAPI(c *gin.Context) {
// 	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1") // Ví dụ API
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	c.JSON(http.StatusOK, gin.H{"data": string(body)})
// }

// // POST /external-api: Kết nối đến API external với POST
// func PostExternalAPI(c *gin.Context) {
// 	data := map[string]string{"title": "foo", "body": "bar"}
// 	jsonData, _ := json.Marshal(data)

// 	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	c.JSON(http.StatusOK, gin.H{"data": string(body)})
// }