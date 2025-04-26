package routes

import (
	"managedata/config"
	"managedata/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	// Ping API
	server.GET(config.AppConfig.Prefix+"/ping", PingResponse)

	server.POST(config.AppConfig.Prefix+"/uploadFile", controller.UploadExcelHandler)
	server.GET(config.AppConfig.Prefix+"/getAllData", controller.ViewImportedList)
	server.GET(config.AppConfig.Prefix+"/getDataById", controller.ViewEmployeeByID)
	server.PATCH(config.AppConfig.Prefix+"/updateById", controller.UpdateEmployee)

}

// Function to Check PING Response
func PingResponse(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "PONG"})
}
