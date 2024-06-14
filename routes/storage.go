package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
	"leanmeal/api/models"
	"leanmeal/api/repositories"
)

type StorageController struct {
	localStorage     string
	EventsRepository repositories.EventRepository
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

	var filePath string
	if err := ctx.Bind(&filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	file, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	name := file.Name()

	ctx.Header("Content-Disposition", "attachment; filename="+name)
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
	var directory string
	if err := ctx.Bind(&directory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	var dirPath string
	if directory == "" {
		dirPath = s.localStorage + directory
	}

	if directory != "" {
		dirPath = directory
	}

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

	ctx.JSON(http.StatusOK, files)
}

func (s *StorageController) listFiles(dirPath string) ([]models.StorageItem, error) {
	items, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var directoryFiles []models.StorageItem

	for _, item := range items {
		wg.Add(1)
		go func(item os.DirEntry) {
			defer wg.Done()

			info, err := item.Info()
			if err != nil {

				return
			}

			if item.IsDir() {
				subDirPath := filepath.Join(dirPath, item.Name())
				count := s.getDirectoryItemCount(subDirPath)
				directorySize, err := s.getDirSize(subDirPath)
				if err != nil {
					return
				}

				mu.Lock()
				directoryFiles = append(directoryFiles, models.StorageItem{
					Size:         directorySize,
					Name:         item.Name(),
					UpdatedAt:    info.ModTime(),
					Type:         1,
					IsDirectory:  true,
					ItemsCount:   count,
					AbsolutePath: subDirPath,
				})
				mu.Unlock()
			} else {
				mu.Lock()
				subDirPath := filepath.Join(dirPath, item.Name())

				directoryFiles = append(directoryFiles, models.StorageItem{
					Size:         info.Size(),
					Name:         item.Name(),
					UpdatedAt:    info.ModTime(),
					Type:         2,
					IsDirectory:  false,
					ItemsCount:   0,
					AbsolutePath: subDirPath,
				})
				mu.Unlock()
			}
		}(item)
	}

	wg.Wait()

	return directoryFiles, nil
}

func (s *StorageController) getDirectoryItemCount(directory string) int {
	items, _ := os.ReadDir(directory)
	var count int
	for range items {

		count++
	}
	return count
}

func (s *StorageController) getDirSize(dirPath string) (int64, error) {
	var size int64

	// Walk through the directory and sum up file sizes
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

func (s *StorageController) Init(r *gin.RouterGroup, m *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("storage")

	// controller.Use(m.Authorize())

	controller.POST("upload", s.upload)
	controller.POST("download", s.download)
	controller.GET("stream/:filename", s.stream)
	controller.POST("path", s.path)
	controller.GET("", s.path)
}
