package routes

import (
	dms "YSKH_DMS/internal/controllers/dms"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("customertype-api", dms.CustomertypeApi)
	r.GET("productcategory-api", dms.ProductcategoryApi)
	r.GET("product-api", dms.ProductApi)
	r.GET("pricelist-api", dms.PriceListApi)
	return r
}
