package models

import "time"

type StorageItem struct {
	Name        string
	Type        int
	Size        int64
	UpdatedAt   time.Time
	ItemsCount  int
	IsDirectory bool
}
