package routes

import (
	"github.com/gin-gonic/gin"
	dms "YSKH_DMS/internal/controllers/dms"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("customertype-api", dms.CustomertypeApi)   
	r.GET("productcategory-api", dms.ProductcategoryApi)  
	r.GET("product-api", dms.ProductApi)   
	r.GET("pricelist-api", dms.PriceListApi)
	r.GET("distunit-api", dms.OrganizationListApiWeb)  
	r.GET("organization-api", dms.OrganizationListApi)    
	return r
}
