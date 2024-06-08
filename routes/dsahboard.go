package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
	"leanmeal/api/models"
	"leanmeal/api/repositories"
)

type DashboardController struct {
	DashboardRepository repositories.DashboardRepository
}

func (d *DashboardController) credentials(ctx *gin.Context) {
	d.DashboardRepository.Storage.Open()
	defer d.DashboardRepository.Storage.Close()

	credentials, err := d.DashboardRepository.Credentials()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, credentials)
}

func (d *DashboardController) devices(ctx *gin.Context) {
	d.DashboardRepository.Storage.Open()
	defer d.DashboardRepository.Storage.Close()

	devices, err := d.DashboardRepository.Devices()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, devices)
}

func (d *DashboardController) storage(ctx *gin.Context) {
	d.DashboardRepository.Storage.Open()
	defer d.DashboardRepository.Storage.Close()

	ctx.JSON(http.StatusOK, models.StorageConsumption{
		Total:          0,
		Used:           0,
		Available:      0,
		TotalUploads:   0,
		TotalDownloads: 0,
		FilesCount:     0,
	})
}

func (d *DashboardController) Init(r *gin.RouterGroup, m *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("dashboard")
	controller.Use(m.Authorize())

	controller.GET("credentials-data", d.credentials)
	controller.GET("devices", d.devices)
	controller.GET("storage", d.storage)

}
