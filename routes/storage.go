package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
)

type StorageController struct {
	localStorage string
}

func (s *StorageController) upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ctx.SaveUploadedFile(file, s.localStorage+file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (s *StorageController) download(ctx *gin.Context) {

	filename := ctx.Param("filename")
	filePath := s.localStorage + filename

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(filePath)

}

func (s *StorageController) stream(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filePath := s.localStorage + filename

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(filePath)
}

func (s *StorageController) path(ctx *gin.Context) {
	directory := ctx.Param("directory")
	dirPath := s.localStorage + directory

	dirInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) || !dirInfo.IsDir() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Directory not found"})
		return
	}

	files, err := s.listFiles(dirPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list directory"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"directory": directory, "files": files})
}

func (s *StorageController) listFiles(dirPath string) ([]string, error) {
	var files []string

	items, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IsDir() {
			files = append(files, item.Name()+"/")
		}
	}

	for _, item := range items {
		if !item.IsDir() {
			files = append(files, item.Name())
		}
	}

	return files, nil
}

func (s *StorageController) Init(r *gin.RouterGroup, m *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("storage")

	controller.Use(m.Authorize())

	controller.POST("upload", s.upload)
	controller.GET("download/:filename", s.download)
	controller.GET("stream/:filename", s.stream)
	controller.GET("path/:directory", s.path)
}
