package dms

import (
	_ "bytes"
	_ "encoding/json"
	_ "io"
	_ "net/http"
	_ "net/url"    
	_ "strings"   
	_ "time"
	"encoding/json"
	vtservice "YSKH_DMS/internal/services/viettel"
	
	 "github.com/gin-gonic/gin"
	_ "YSKH_DMS/internal/models"
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
	queryApi := map[string]any{
			"status":     []string{"ACT", "PEN"},
			"createDate": "2025-05-16",
		}
	urlapi := "https://app.vietteldms.com/openapi/v1/GetListDist"
	urlApiBuild, _ :=  vtservice.BuildViettelDistURL(urlapi,queryApi)
	repData, _ :=  vtservice.ViettelGet(urlApiBuild)
	// c.JSON(200, gin.H{"data": string(repData)})
	var result any
	json.Unmarshal(repData, &result)
	c.JSON(200, gin.H{
		"success": true,
		"data":    result, // ← object hoặc array thật
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