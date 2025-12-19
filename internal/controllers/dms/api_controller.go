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
func DemoApi(c *gin.Context) {

	repData, err := customertypeService.SaveBatchDmsViettel();
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    repData,
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