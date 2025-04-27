package routes

import (
	"managedata/config"
	"managedata/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	server.GET(config.AppConfig.Prefix+"/ping", PingResponse)

	server.POST(config.AppConfig.Prefix+"/uploadFile", controller.UploadExcelHandler)
	server.GET(config.AppConfig.Prefix+"/getAllData", controller.GetAllDataHandler)
	server.GET(config.AppConfig.Prefix+"/getDataById", controller.GetDataByIdHandler)
	server.PATCH(config.AppConfig.Prefix+"/updateById", controller.UpdataDataByIdHandler)
	server.DELETE(config.AppConfig.Prefix+"/deleteById", controller.DeleteEmployeeHandler)

}

// Function to Check PING Response
func PingResponse(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": "PONG"})
}
