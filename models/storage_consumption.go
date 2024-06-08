package models

type StorageConsumption struct {
	Total          int
	Used           int
	Available      int
	TotalUploads   int
	TotalDownloads int
	FilesCount     int
}
