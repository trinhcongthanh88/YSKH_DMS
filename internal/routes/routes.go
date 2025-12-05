package routes

import (
	"github.com/gin-gonic/gin"
	dms "YSKH_DMS/internal/controllers/dms"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("demo-api", dms.DemoApi)   
	return r
}